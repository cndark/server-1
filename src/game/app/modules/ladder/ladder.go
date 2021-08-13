package ladder

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	. "fw/src/core/math"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/robot"
	"fw/src/game/app/modules/svrgrp"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"math/rand"
	"sort"
	"time"
)

// ============================================================================

const (
	c_Stage_Close   = 0
	c_Stage_Init    = 1
	c_Stage_Robot1  = 2
	c_Stage_Robot2  = 3
	c_Stage_Prepare = 4
	c_Stage_Start   = 5
	c_Stage_Reward  = 6
)

const (
	c_FightLockSec = 30
)

// ============================================================================

var (
	G *svrgrp.Group // 跨服组

	g_stage int       // 当前阶段
	g_t0    time.Time // 本期版本
	g_t2    time.Time // 当前阶段结束时间 (下个阶段开始时间)

	g_svrlen  int32          // 本期服务器组有效长度
	g_isin    bool           // 本服是否参与
	g_rwdsent map[int32]bool // 是否发奖. [svrid]

	ladder_rank  []*ladder_plr_t          // rank
	ladder_index map[string]*ladder_plr_t // index by id

	ladder_ft_locks map[string]time.Time // fight locks
)

// ============================================================================

type ladder_plr_t struct {
	Info    *msg.PlayerSimpleInfo
	Rank    int32
	IsRobot bool
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_SvrGrpReady, func(...interface{}) {
		G = svrgrp.GetGroupCross()
	})

	evtmgr.On(gconst.Evt_GsPull_Ladder_AddPlrRequest, func(args ...interface{}) {
		oarg := args[1].([]byte)
		ret := args[2].(func(int32, interface{}))

		// check args
		var ctx *ladder_plr_t
		err := utils.UnmarshalArg(oarg, &ctx)
		if err != nil {
			log.Error("unmarshal ladder-add-plr-request ctx failed:", err)
			ret(Err.Failed, nil)
			return
		}

		// check if exists
		e := ladder_index[ctx.Info.Id]
		if e != nil {
			ret(Err.OK, nil)
			return
		}

		// add
		ladder_index[ctx.Info.Id] = ctx
		ladder_rank = append(ladder_rank, ctx)
		ctx.Rank = int32(len(ladder_rank))

		// broadcast, first
		G.PushAll(gconst.Evt_GsPush_Ladder_AddPlr, nil, ctx, true)

		// save to db
		save_addplr(ctx)

		// reply
		ret(Err.OK, nil)
	})

	evtmgr.On(gconst.Evt_GsPush_Ladder_AddPlr, func(args ...interface{}) {
		oarg := args[1].([]byte)

		if !g_isin {
			return
		}

		// check args
		var ctx *ladder_plr_t
		err := utils.UnmarshalArg(oarg, &ctx)
		if err != nil {
			log.Error("unmarshal ladder-add-plr ctx failed:", err)
			return
		}

		// add
		ladder_index[ctx.Info.Id] = ctx
		ladder_rank = append(ladder_rank, ctx)
	})

	evtmgr.On(gconst.Evt_GsPull_Ladder_FightLock, func(args ...interface{}) {
		ids := args[0].([]string)
		ret := args[2].(func(int32, interface{}))

		// check
		if ladder_ft_locks == nil {
			ladder_ft_locks = make(map[string]time.Time)
		}

		now := time.Now()
		for _, id := range ids {
			if now.Sub(ladder_ft_locks[id]).Seconds() < c_FightLockSec {
				ret(Err.Failed, nil)
				return
			}
		}

		// lock
		for _, id := range ids {
			ladder_ft_locks[id] = now
		}

		ret(Err.OK, nil)
	})

	evtmgr.On(gconst.Evt_GsPush_Ladder_FightUnlock, func(args ...interface{}) {
		ids := args[0].([]string)

		for _, id := range ids {
			delete(ladder_ft_locks, id)
		}
	})

	evtmgr.On(gconst.Evt_GsPush_Ladder_SyncPos, func(args ...interface{}) {
		sarg := args[0].([]string)

		rk1 := core.Atoi32(sarg[0])
		rk2 := core.Atoi32(sarg[1])

		// check
		L := int32(len(ladder_rank))
		if rk1 < 1 || rk1 > L || rk2 < 1 || rk2 > L {
			return
		}

		// change pos
		e1 := ladder_rank[rk1-1]
		e2 := ladder_rank[rk2-1]

		ladder_rank[rk1-1], ladder_rank[rk2-1] = e2, e1
		e1.Rank = rk2
		e2.Rank = rk1
	})
}

// ============================================================================

// db format:
/*
	_id: grpid
	t0: t0
	svrlen: n

	rwdsent: {
		svrid: true,
	}

	robot: {
		done: false,
		svrs: {
			svrid: true,
		}
	}

	plrs: {
		id: ladder_plr_t
	}
*/

func on_stage(stg string, t0, t1, t2 time.Time, f func(bool)) {
	// set ts
	g_t0 = t0
	g_t2 = t2

	if stg == "init" {
		g_stage = c_Stage_Init

		ladder_index = nil
		ladder_rank = nil
		ladder_ft_locks = make(map[string]time.Time)

		if G.IsMaster() {
			async.Push(func() {
				err := dbmgr.DBCross.Upsert(
					dbmgr.C_tabname_ladder,
					G.Id,
					db.M{
						"$set": db.M{
							"t0":         t0,
							"svrlen":     G.AvailLen,
							"rwdsent":    db.M{},
							"robot.done": false,
						},
						"$unset": db.M{
							"plrs":       1,
							"robot.svrs": 1,
						},
					},
				)
				if err != nil {
					log.Error("init ladder db failed:", err)
				}
			})
		}
	} else {
		core.Go(func() {
			var obj struct {
				T0     time.Time
				SvrLen int32
			}
			err := dbmgr.DBCross.GetProjection(
				dbmgr.C_tabname_ladder,
				G.Id,
				db.M{"t0": 1, "svrlen": 1},
				&obj,
			)
			loop.Push(func() {
				// only those available at 'init' phase are valid servers
				b := false
				for i := int32(0); i < obj.SvrLen; i++ {
					if G.Svrs[i] == config.CurGame.Id {
						b = true
						break
					}
				}

				g_svrlen = obj.SvrLen
				g_isin = err == nil && obj.T0.Equal(t0) && b

				f(g_isin)
			})
		})
	}
}

// ============================================================================

func stage_robot1() {
	set_stage(c_Stage_Robot1)

	core.Go(func() {
		var obj struct {
			Robot struct {
				Svrs map[int32]bool
			}
		}

		err := dbmgr.DBCross.GetProjection(
			dbmgr.C_tabname_ladder,
			G.Id,
			db.M{"robot": 1},
			&obj,
		)
		if err != nil {
			log.Error("ladder robot1 failed:", err)
		}

		loop.Push(func() {
			if obj.Robot.Svrs[config.CurGame.Id] {
				return
			}

			// contribute local robots
			conf := gamedata.ConfGlobalPublic.Query(1)
			if conf == nil {
				return
			}

			arr := make([]*ladder_plr_t, 0, 100)
			for _, v := range conf.LadderRobot {
				for _, v2 := range robot.RobotMgr.GetByLevel(v.Lv, v.N/g_svrlen+1) {
					arr = append(arr, &ladder_plr_t{
						Info: &msg.PlayerSimpleInfo{
							Id:       v2.Id,
							Name:     v2.Name,
							Lv:       v2.Lv,
							Exp:      0,
							Head:     v2.Head,
							HFrame:   v2.HFrame,
							Vip:      0,
							SvrId:    config.CurGame.Id,
							AtkPwr:   v2.AtkPwr,
							GName:    "",
							ShowHero: v2.GetShowHero(),
						},
						IsRobot: true,
					})
				}
			}

			// save
			async.Push(func() {
				err := dbmgr.DBCross.Upsert(
					dbmgr.C_tabname_ladder,
					G.Id,
					db.M{
						"$push": db.M{
							"plrs": db.M{
								"$each": arr,
							},
						},
						"$set": db.M{
							fmt.Sprintf("robot.svrs.%d", config.CurGame.Id): true,
						},
					},
				)
				if err != nil {
					log.Error("contribute local robots failed:", err)
				}
			})
		})
	})
}

func stage_robot2() {
	set_stage(c_Stage_Robot2)

	if !G.IsMaster() {
		return
	}

	core.Go(func() {
		var obj struct {
			Plrs []*ladder_plr_t
		}

		err := dbmgr.DBCross.GetProjectionByCond(
			dbmgr.C_tabname_ladder,
			db.M{
				"_id":        G.Id,
				"robot.done": false,
			},
			db.M{"plrs": 1},
			&obj,
		)
		if db.IsNotFound(err) {
			return
		} else if err != nil {
			log.Error("ladder robot2 failed:", err)
		}

		loop.Push(func() {
			// shuffle
			rand.Shuffle(len(obj.Plrs), func(i, j int) {
				obj.Plrs[i], obj.Plrs[j] = obj.Plrs[j], obj.Plrs[i]
			})

			// sort
			sort.Slice(obj.Plrs, func(i, j int) bool {
				return obj.Plrs[i].Info.Lv > obj.Plrs[j].Info.Lv
			})

			// assign rank
			for i, v := range obj.Plrs {
				v.Rank = int32(i) + 1
			}

			// convert to map
			m := make(map[string]*ladder_plr_t)
			for _, v := range obj.Plrs {
				m[v.Info.Id] = v
			}

			// save
			async.Push(func() {
				err := dbmgr.DBCross.Upsert(
					dbmgr.C_tabname_ladder,
					G.Id,
					db.M{
						"$set": db.M{
							"plrs":       m,
							"robot.done": true,
						},
					},
				)
				if err != nil {
					log.Error("organizing robots failed:", err)
				}
			})
		})
	})
}

func stage_prepare() {
	set_stage(c_Stage_Prepare)
}

func stage_start() {
	load_ladder_plrs(func() {
		set_stage(c_Stage_Start)
	})
}

func stage_reward() {
	load_ladder_plrs(func() {
		set_stage(c_Stage_Reward)

		send_rewards_mails()
	})
}

func stage_close() {
	set_stage(c_Stage_Close)
}

// ============================================================================

func set_stage(v int) {
	g_stage = v
	utils.BroadcastPlayers(&msg.GS_LadderStageChange{
		Stage: int32(v),
		Ts2:   g_t2.Unix(),
	})
}

func load_ladder_plrs(f func()) {
	core.Go(func() {
		var obj struct {
			RwdSent map[int32]bool
			Plrs    map[string]*ladder_plr_t
		}

		err := dbmgr.DBCross.GetObject(
			dbmgr.C_tabname_ladder,
			G.Id,
			&obj,
		)
		if err != nil {
			log.Error("loading ladder players failed:", err)
		}

		// make rank array
		arr := make([]*ladder_plr_t, 0, len(obj.Plrs))
		for _, v := range obj.Plrs {
			arr = append(arr, v)
		}

		// sort by rank
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].Rank < arr[j].Rank
		})

		loop.Push(func() {
			g_rwdsent = obj.RwdSent
			ladder_rank = arr
			ladder_index = obj.Plrs

			f()
		})
	})
}

func add_ladder_plr(lp *ladder_plr_t, f func(ec int32)) {
	utils.GsPull(
		G.Svrs[0], // master
		gconst.Evt_GsPull_Ladder_AddPlrRequest,
		nil, lp, nil,
		func(ec int32, r interface{}) {
			f(ec)
		},
	)
}

func fight_lock(ids []string, f func(ec int32)) {
	utils.GsPull(
		G.Svrs[0], // master
		gconst.Evt_GsPull_Ladder_FightLock,
		ids, nil, nil,
		func(ec int32, r interface{}) {
			f(ec)
		},
	)
}

func fight_unlock(ids []string) {
	// delayed
	loop.SetTimeout(time.Now().Add(time.Millisecond*1000), func() {
		G.Push2Master(gconst.Evt_GsPush_Ladder_FightUnlock, ids, nil)
	})
}

func fight_do(ctx *fight_ctx_t, ret func(int32, interface{})) {
	// find attacker
	e1 := ladder_index[ctx.T1.Player.Id]
	if e1 == nil {
		ret(Err.Ladder_TargetNotFound, nil)
		return
	}

	// find target
	e2 := ladder_index[ctx.TarId]
	if e2 == nil {
		ret(Err.Ladder_TargetNotFound, nil)
		return
	}

	// fight
	bi := &msg.BattleInput{
		T1: ctx.T1,
		T2: e2.to_team(),
		Args: map[string]string{
			"Module": "LADDER",
			"Round":  "0",
		},
	}

	battle.Fight(bi, func(r *msg.BattleResult) {
		// check
		if r == nil {
			ret(Err.Failed, nil)
			return
		}

		// replay
		replay := &msg.BattleReplay{
			Bi: bi,
			Br: r,
		}

		// should change pos ?
		b := r.Winner == 1 && e1.Rank > e2.Rank

		// ret
		ret(Err.OK, &fight_res_t{
			Replay: replay,
			ChgPos: b,
		})

		// add def replay
		e2.add_def_replay(replay, e1.Rank, e2.Rank, b)
	})
}

func sync_pos(rk1, rk2 int32) {
	G.PushAll(gconst.Evt_GsPush_Ladder_SyncPos, []string{core.I32toa(rk1), core.I32toa(rk2)}, nil)

	// save to db
	save_chgpos(rk1, rk2)
}

func save_addplr(p *ladder_plr_t) {
	async.Push(func() {
		err := dbmgr.DBCross.Upsert(
			dbmgr.C_tabname_ladder,
			G.Id,
			db.M{
				"$set": db.M{
					fmt.Sprintf("plrs.%s", p.Info.Id): p,
				},
			},
		)
		if err != nil {
			log.Error("save ladder add-plr failed:", err)
		}
	})
}

func save_chgpos(rk1, rk2 int32) {
	// check
	L := int32(len(ladder_rank))
	if rk1 < 1 || rk1 > L || rk2 < 1 || rk2 > L {
		return
	}

	// find ids
	id1 := ladder_rank[rk1-1].Info.Id
	id2 := ladder_rank[rk2-1].Info.Id

	// save
	async.Push(func() {
		err := dbmgr.DBCross.Upsert(
			dbmgr.C_tabname_ladder,
			G.Id,
			db.M{
				"$set": db.M{
					fmt.Sprintf("plrs.%s.rank", id1): rk2,
					fmt.Sprintf("plrs.%s.rank", id2): rk1,
				},
			},
		)
		if err != nil {
			log.Error("save ladder chg-pos failed:", err)
		}
	})
}

func send_rewards_mails() {
	// already sent ?
	if g_rwdsent[config.CurGame.Id] {
		return
	}

	// reward
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	L := int32(len(ladder_rank))

	for _, v := range gamedata.ConfLadder.Items() {
		a := MaxInt32(1, v.RankRange[0].Up)
		b := MinInt32(L, v.RankRange[0].Down)

		for i := a; i <= b; i++ {
			e := ladder_rank[i-1]

			// fuck robot
			if e.IsRobot {
				continue
			}

			// only local
			if e.Info.SvrId != config.CurGame.Id {
				continue
			}

			// find player
			iplr := utils.LoadPlayer(e.Info.Id)
			if iplr == nil {
				continue
			}

			// send
			ml := mail.New(iplr).SetKey(conf.LadderRewardMailId)
			for _, v2 := range v.Reward {
				ml.AddAttachment(v2.Id, float64(v2.N))
			}
			ml.AddDictInt32("rank", i)
			ml.Send()
		}
	}

	// mark as sent
	g_rwdsent[config.CurGame.Id] = true

	async.Push(func() {
		err := dbmgr.DBCross.Upsert(
			dbmgr.C_tabname_ladder,
			G.Id,
			db.M{
				"$set": db.M{
					fmt.Sprintf("rwdsent.%d", config.CurGame.Id): true,
				},
			},
		)
		if err != nil {
			log.Error("mark ladder reward mail sent status failed:", err)
		}
	})
}

// ============================================================================

// MUST be local svr entities
func (self *ladder_plr_t) to_team() *msg.BattleTeam {
	if self.IsRobot {
		rb := robot.RobotMgr.FindRobot(self.Info.Id)
		if rb == nil {
			return nil
		} else {
			return rb.ToMsg_BattleTeam(rb.Team)
		}
	} else {
		iplr := utils.LoadPlayer(self.Info.Id)
		if iplr == nil {
			return nil
		} else {
			plr := iplr.(IPlayer)
			return plr.ToMsg_BattleTeam(plr.GetTeam(gconst.TeamType_Dfd))
		}
	}
}

func (self *ladder_plr_t) add_def_replay(rp *msg.BattleReplay, rk1, rk2 int32, chgpos bool) {
	if !self.IsRobot {
		iplr := utils.LoadPlayer(self.Info.Id)
		if iplr != nil {
			iplr.(IPlayer).GetLadder().add_def_replay(rp, rk1, rk2, chgpos)
		}
	}
}
