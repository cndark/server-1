package ladder

import (
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math"
	"time"
)

// ============================================================================

const (
	c_DefendReplayMax = 10
)

// ============================================================================

type Ladder struct {
	Replays []*def_replay_t

	plr IPlayer
}

type fight_ctx_t struct {
	T1    *msg.BattleTeam
	TarId string
}

type fight_res_t struct {
	Replay *msg.BattleReplay
	ChgPos bool
}

type def_replay_t struct {
	ReplayId string
	Attacker *msg.PlayerSimpleInfo
	Winner   int32
	Ts       int64
	RkFrom   int32
	RkTo     int32
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GsPull_Ladder_Fight, func(args ...interface{}) {
		oarg := args[1].([]byte)
		ret := args[2].(func(int32, interface{}))

		// check args
		var ctx *fight_ctx_t
		err := utils.UnmarshalArg(oarg, &ctx)
		if err != nil {
			log.Error("unmarshal ladder-fight ctx failed:", err)
			ret(Err.Failed, nil)
			return
		}

		// fight
		fight_do(ctx, ret)
	})
}

// ============================================================================

func NewLadder() *Ladder {
	return &Ladder{}
}

func (self *Ladder) Init(plr IPlayer) {
	self.plr = plr
}

// ============================================================================

func (self *Ladder) Match(f func(ec int32, ret []*msg.LadderPlayerInfo)) {
	// check stage
	if g_stage != c_Stage_Start {
		f(Err.Failed, nil)
		return
	}

	// make sure we're joined
	self.join(func(ec int32) {
		// check err
		if ec != Err.OK {
			f(ec, nil)
			return
		}

		// get entry
		e := ladder_index[self.plr.GetId()]
		if e == nil {
			f(Err.Failed, nil)
			return
		}

		// conf
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			f(Err.Failed, nil)
			return
		}

		// match
		a := map[int32]bool{1: true, 2: true, 3: true}

		// use coeff ?
		L := int32(len(ladder_rank))
		if e.Rank > 50 {
			// scan match coeff
			for _, c := range conf.LadderMatch {
				i := int32(math.Floor(float64(e.Rank) * c))
				if i < 4 {
					continue
				}
				if i > L {
					break
				}

				a[i] = true

				if int32(len(a)) >= conf.LadderMatchNum {
					break
				}
			}
		} else if e.Rank < conf.LadderMatchNum {
			for i := int32(4); i <= conf.LadderMatchNum; i++ {
				if i > L {
					break
				}

				a[i] = true
			}
		} else {
			// index diff should be 1
			for i := e.Rank - (conf.LadderMatchNum - 5); ; i++ {
				if i < 4 {
					continue
				}
				if i > L {
					break
				}

				a[i] = true

				if int32(len(a)) >= conf.LadderMatchNum {
					break
				}
			}
		}

		// ok. now 'a' holds opponents indexes
		ret := make([]*msg.LadderPlayerInfo, 0, conf.LadderMatchNum)
		for i := range a {
			if i-1 >= L {
				continue
			}

			e := ladder_rank[i-1]
			ret = append(ret, &msg.LadderPlayerInfo{
				Info: e.Info,
				Rank: e.Rank,
			})

		}

		f(Err.OK, ret)
	})
}

func (self *Ladder) join(f func(ec int32)) {
	// already joined ?
	e := ladder_index[self.plr.GetId()]
	if e != nil {
		f(Err.OK)
		return
	}

	// check def-team
	if !self.plr.IsSetTeam(gconst.TeamType_Dfd) {
		f(Err.Common_NotSetTeam)
		return
	}

	// make join request
	lp := &ladder_plr_t{
		Info: self.plr.ToMsg_SimpleInfo(),
	}
	add_ladder_plr(lp, f)
}

func (self *Ladder) Fight(tf *msg.TeamFormation, tarid string, rk int32, f func(ec int32, replay *msg.BattleReplay, rwd *msg.Rewards)) {
	// check stage
	if g_stage != c_Stage_Start {
		f(Err.Failed, nil, nil)
		return
	}

	// check attacker
	e1 := ladder_index[self.plr.GetId()]
	if e1 == nil {
		f(Err.Ladder_TargetNotFound, nil, nil)
		return
	}

	// check target
	e2 := ladder_index[tarid]
	if e2 == nil {
		f(Err.Ladder_TargetNotFound, nil, nil)
		return
	}
	if e2.Rank != rk {
		f(Err.Ladder_TargetChanged, nil, nil)
		return
	}

	// check team
	if !self.plr.IsTeamFormationValid(tf) {
		f(Err.Plr_TeamInvalid, nil, nil)
		return
	}

	// cost
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_LadderFight)
	op.Dec(20011, 1)
	if ec := op.CheckEnough(); ec != Err.OK {
		f(ec, nil, nil)
		return
	}

	// reward
	for id, n := range utils.BattleReward(self.plr.GetLevel()) {
		op.Inc(id, n)
	}

	// acquire fight lock
	ids := []string{self.plr.GetId(), tarid}

	fight_lock(ids, func(ec int32) {
		// check err
		if ec != Err.OK {
			f(Err.Ladder_InFight, nil, nil)
			return
		}

		// check target again
		if e2.Rank != rk {
			fight_unlock(ids)
			f(Err.Ladder_TargetChanged, nil, nil)
			return
		}

		// go to target svr to fight
		utils.GsPull(
			e2.Info.SvrId,
			gconst.Evt_GsPull_Ladder_Fight,
			nil, &fight_ctx_t{
				T1:    self.plr.ToMsg_BattleTeam(tf, true),
				TarId: tarid,
			}, &fight_res_t{},
			func(ec int32, r interface{}) {
				defer fight_unlock(ids)

				if ec != Err.OK {
					f(ec, nil, nil)
					return
				}

				// ok
				rwd := op.Apply().ToMsg()
				res := r.(*fight_res_t)

				// change pos if should
				if res.ChgPos {
					sync_pos(e1.Rank, e2.Rank)
				}

				// callback
				f(Err.OK, res.Replay, rwd)

				// evt
				evtmgr.Fire(gconst.Evt_LadderFight, self.plr)
			},
		)
	})
}

func (self *Ladder) add_def_replay(rp *msg.BattleReplay, rk1, rk2 int32, chgpos bool) {
	// save
	rpid := battle.ReplaySave(dbmgr.DBGame, dbmgr.C_tabname_replay_ladder, rp)
	bplr := rp.Bi.T1.Player

	// calc rk2_
	rk2_ := rk2
	if chgpos {
		rk2_ = rk1
	}

	e := &def_replay_t{
		ReplayId: rpid,
		Attacker: &msg.PlayerSimpleInfo{
			Id:     bplr.Id,
			Name:   bplr.Name,
			Lv:     bplr.Lv,
			Exp:    0,
			Head:   bplr.Head,
			HFrame: bplr.HFrame,
			Vip:    bplr.Vip,
			SvrId:  bplr.SvrId,
			AtkPwr: bplr.AtkPwr,
			GName:  "",
		},
		Winner: rp.Br.Winner,
		Ts:     time.Now().Unix(),
		RkFrom: rk2,
		RkTo:   rk2_,
	}

	self.Replays = append(self.Replays, e)

	L := len(self.Replays)
	if L > c_DefendReplayMax {
		battle.ReplayDel(dbmgr.DBGame, dbmgr.C_tabname_replay_ladder, self.Replays[0].ReplayId)
		self.Replays = self.Replays[1:]
	}
}

func (self *Ladder) GetLadderRank() (ret []*msg.LadderPlayerInfo) {
	// check stage
	if g_stage < c_Stage_Start {
		return
	}

	// get
	L := len(ladder_rank)
	if L > 100 {
		L = 100
	}

	ret = make([]*msg.LadderPlayerInfo, 0, L)
	for _, v := range ladder_rank[:L] {
		ret = append(ret, &msg.LadderPlayerInfo{
			Info: v.Info,
			Rank: v.Rank,
		})
	}

	return
}

func (self *Ladder) GetSelfRank() int32 {
	if g_stage >= c_Stage_Start {
		e := ladder_index[self.plr.GetId()]
		if e != nil {
			return e.Rank
		}
	}

	return 0
}

func (self *Ladder) GetReplayList() (ret []*msg.LadderReplayRec) {
	ret = make([]*msg.LadderReplayRec, 0, c_DefendReplayMax)
	for _, v := range self.Replays {
		ret = append(ret, (*msg.LadderReplayRec)(v))
	}

	return
}

func (self *Ladder) GetReplay(id string, f func(*msg.BattleReplay)) {
	battle.ReplayGet(dbmgr.DBGame, dbmgr.C_tabname_replay_ladder, id, func(rp *msg.BattleReplay) {
		f(rp)
	})
}

// ============================================================================

func (self *Ladder) ToMsg() *msg.LadderData {
	return &msg.LadderData{
		Stage: int32(g_stage),
		Ts2:   g_t2.Unix(),
	}
}
