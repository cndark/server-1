package gwar

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/mdata"
	"fw/src/game/app/modules/svrgrp"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"sort"
	"time"
)

// ============================================================================

const (
	c_Stage_Close   = 0
	c_Stage_Init    = 1
	c_Stage_Prepare = 2
	c_Stage_Enroll  = 3
	c_Stage_Match   = 4
	c_Stage_Reward  = 5
)

// ============================================================================

var (
	G *svrgrp.Group // 跨服组

	localdata *localdata_t

	g_stage int       // 当前阶段
	g_t0    time.Time // 本期版本
	g_t2    time.Time // 当前阶段结束时间 (下个阶段开始时间)

	enroll_list  []*enroll_base_info_t          // sorted by ap. used for matching
	enroll_rank  []*enroll_base_info_t          // sorted by jf
	enroll_index map[string]*enroll_base_info_t // index
)

// ============================================================================

// 本服数据
type localdata_t struct {
	Ver      time.Time              // ver
	Glds     map[string]*gld_data_t // [gid]
	Enrolled bool                   // enrolled flag
}

type gld_data_t struct {
	Id          string
	G1Plrs      map[string]*g1_plrinfo_t // [plrid]
	G2          *enroll_gld_t
	G2Plrs      map[string]*g2_plrinfo_t // [plrid]
	G2DeadCount int
	Hist        map[string]bool // [gid]
}

type g1_plrinfo_t struct {
	Cnt    int32
	Jf     int32
	RwdGld bool // rwd taken ?
	RwdPlr bool // rwd taken ?
}

type g2_plrinfo_t struct {
	locked bool
	Done   bool
	Val    int32
}

type enroll_gld_t struct {
	Base *enroll_base_info_t
	Mbs  map[string]*msg.BattleTeam // [plrid]
}

type enroll_base_info_t struct {
	Id     string
	Name   string
	Icon   int32
	Lv     int32
	N      int32
	ApMean int32
	ApMax  int32
	SvrId  int32
	Jf     int32
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_SvrGrpReady, func(...interface{}) {
		G = svrgrp.GetGroupCross()
	})

	evtmgr.On(gconst.Evt_GsPush_GWarGldJf, func(args ...interface{}) {
		sarg := args[0].([]string)

		gid := sarg[0]
		jf_add := core.Atoi32(sarg[1])

		e := enroll_index[gid]
		if e != nil {
			e.Jf += jf_add

			// sort rank
			sort.Slice(enroll_rank, func(i, j int) bool {
				return enroll_rank[i].Jf > enroll_rank[j].Jf
			})
		}
	})
}

// ============================================================================

func new_local_data() interface{} {
	return &localdata_t{
		Glds: make(map[string]*gld_data_t),
	}
}

func on_local_data_loaded() {
	localdata = mdata.Get(NAME).(*localdata_t)
}

// ============================================================================

// db format:
/*
	_id: grpid
	t0: t0

	enrolls: []*enroll_gld_t
*/

func on_stage(stg string, t0, t1, t2 time.Time, f func(bool)) {
	// data reset
	if !localdata.Ver.Equal(t0) {
		localdata.Ver = t0
		localdata.reset_data()
	}

	// set ts
	g_t0 = t0
	g_t2 = t2

	if stg == "init" {
		g_stage = c_Stage_Init

		if G.IsMaster() {
			async.Push(func() {
				err := dbmgr.DBCross.Upsert(
					dbmgr.C_tabname_gwar,
					G.Id,
					db.M{
						"$set": db.M{
							"t0": t0,
						},
						"$unset": db.M{
							"enrolls": 1,
						},
					},
				)
				if err != nil {
					log.Error("init gwar db failed:", err)
				}
			})
		}
	} else {
		core.Go(func() {
			var obj struct {
				T0 time.Time ``
			}
			err := dbmgr.DBCross.GetProjection(
				dbmgr.C_tabname_gwar,
				G.Id,
				db.M{"t0": 1},
				&obj,
			)
			loop.Push(func() {
				f(err == nil && obj.T0.Equal(t0))
			})
		})
	}
}

// ============================================================================

func stage_prepare() {
	set_stage(c_Stage_Prepare)
}

func stage_enroll() {
	set_stage(c_Stage_Enroll)

	if localdata.Enrolled {
		return
	}

	arr := make([]*enroll_gld_t, 0, 50)
	for _, gld := range guild.GuildMgr.Array_Guilds() {
		if gld.Lv < 4 {
			continue
		}

		mbs := make(map[string]*msg.BattleTeam)
		ap_mean := int64(0)
		ap_max := int32(0)
		for _, mb := range gld.Members {
			iplr := utils.LoadPlayer(mb.Id)
			if iplr != nil {
				plr := iplr.(IPlayer)

				if plr.IsModuleOpen(gconst.ModuleId_Arena) &&
					plr.IsSetTeam(gconst.TeamType_Dfd) {

					t := plr.ToMsg_BattleTeam(plr.GetTeam(gconst.TeamType_Dfd), true)
					mbs[plr.GetId()] = t

					ap := t.Player.AtkPwr
					ap_mean += int64(ap)
					if ap > ap_max {
						ap_max = ap
					}
				}
			}
		}

		L := len(mbs)
		if L == 0 {
			continue
		}

		ap_mean /= int64(L)

		// ok
		localdata.Glds[gld.Id] = &gld_data_t{
			Id:          gld.Id,
			G1Plrs:      make(map[string]*g1_plrinfo_t),
			G2:          nil,
			G2Plrs:      nil,
			G2DeadCount: 0,
			Hist:        make(map[string]bool),
		}

		arr = append(arr, &enroll_gld_t{
			Base: &enroll_base_info_t{
				Id:     gld.Id,
				Name:   gld.Name,
				Icon:   gld.Icon,
				Lv:     gld.Lv,
				N:      int32(L),
				ApMean: int32(ap_mean),
				ApMax:  ap_max,
				SvrId:  config.CurGame.Id,
				Jf:     0,
			},
			Mbs: mbs,
		})
	}

	localdata.Enrolled = true
	async.Push(func() {
		err := dbmgr.DBCross.Upsert(
			dbmgr.C_tabname_gwar,
			G.Id,
			db.M{
				"$push": db.M{
					"enrolls": db.M{
						"$each": arr,
					},
				},
			},
		)
		if err != nil {
			log.Error("gwar saving enrolls failed:", err)
		}
	})
}

func stage_match() {
	load_enroll_list(func() {
		set_stage(c_Stage_Match)

		// match for every local guild
		i := 0
		for _, gd := range localdata.Glds {
			do_match(i*200, gd)
			i++
		}
	})
}

func stage_reward() {
	load_enroll_list(func() {
		set_stage(c_Stage_Reward)

		send_gld_rewards_mail()
		send_plr_rewards_mail()
	})
}

func stage_close() {
	set_stage(c_Stage_Close)
}

// ============================================================================

func set_stage(v int) {
	g_stage = v
	utils.BroadcastPlayers(&msg.GS_GWarStageChange{
		Stage: int32(v),
		Ts2:   g_t2.Unix(),
	})
}

func load_enroll_list(f func()) {
	core.Go(func() {
		var obj struct {
			Enrolls []*enroll_gld_t
		}

		err := dbmgr.DBCross.GetProjection(
			dbmgr.C_tabname_gwar,
			G.Id,
			db.M{"enrolls.base": 1},
			&obj,
		)
		if err != nil && !db.IsNotFound(err) {
			log.Error("loading enroll list failed:", err)
		}

		// make list
		lst := make([]*enroll_base_info_t, 0, len(obj.Enrolls))
		for _, v := range obj.Enrolls {
			lst = append(lst, v.Base)
		}

		// sort list
		sort.Slice(lst, func(i, j int) bool {
			return lst[i].ApMean > lst[j].ApMean
		})

		// make rank
		rk := make([]*enroll_base_info_t, len(lst))
		copy(rk, lst)

		// sort rank
		sort.Slice(rk, func(i, j int) bool {
			return rk[i].Jf > rk[j].Jf
		})

		// make index
		index := make(map[string]*enroll_base_info_t)
		for _, v := range lst {
			index[v.Id] = v
		}

		loop.Push(func() {
			enroll_list = lst
			enroll_rank = rk
			enroll_index = index
			f()
		})
	})
}

func load_enroll_gld(gid string, f func(*enroll_gld_t)) {
	core.Go(func() {
		var obj struct {
			Enrolls []*enroll_gld_t
		}

		err := dbmgr.DBCross.GetProjectionByCond(
			dbmgr.C_tabname_gwar,
			db.M{
				"_id":             G.Id,
				"enrolls.base.id": gid,
			},
			db.M{"enrolls.$": 1},
			&obj,
		)
		if err != nil && !db.IsNotFound(err) {
			log.Error("loading enroll gld failed:", err)
		}

		var eg *enroll_gld_t

		if len(obj.Enrolls) > 0 {
			eg = obj.Enrolls[0]
		}

		loop.Push(func() {
			f(eg)
		})
	})
}

func do_match(delay int, gd *gld_data_t) {
	// get current ap
	gld := guild.GuildMgr.FindGuild(gd.Id)
	if gld == nil {
		return
	}
	ap := gld.GetAtkPwr()

	// find closest by ap
	L := len(enroll_list)
	if L < 2 {
		return
	}

	start := sort.Search(L, func(i int) bool {
		return enroll_list[i].ApMean <= ap
	})
	if start >= L {
		start = L - 1
	}

	idx := -1
	for d := 0; d < L; d++ {
		i := start + d
		if i < L && enroll_list[i].Id != gd.Id && !gd.Hist[enroll_list[i].Id] {
			idx = i
			break
		}

		i = start - d
		if i >= 0 && enroll_list[i].Id != gd.Id && !gd.Hist[enroll_list[i].Id] {
			idx = i
			break
		}
	}

	if idx == -1 {
		for j := 0; j < 100; j++ {
			idx = core.RandInt(0, L-1)
			if enroll_list[idx].Id != gd.Id {
				break
			}
		}
	}

	// matched
	gid2 := enroll_list[idx].Id
	gd.Hist[gid2] = true

	// load
	core.Go(func() {
		time.Sleep(time.Millisecond * time.Duration(delay))
		load_enroll_gld(gid2, func(g2 *enroll_gld_t) {
			gd.G2 = g2
			gd.G2Plrs = make(map[string]*g2_plrinfo_t)
			gd.G2DeadCount = 0

			gld.Broadcast(&msg.GS_GWarNewG2{})
		})
	})
}

func on_gld_win(gd *gld_data_t) {
	// calc jf
	jf_add := gd.G2.Base.ApMax/100 + 10

	// update jf to db
	async.Push(func() {
		err := dbmgr.DBCross.UpsertByCond(
			dbmgr.C_tabname_gwar,
			db.M{
				"_id":             G.Id,
				"enrolls.base.id": gd.Id,
			},
			db.M{
				"$inc": db.M{
					"enrolls.$.base.jf": jf_add,
				},
			},
		)
		if err != nil {
			log.Error("update guild jf failed:", err)
		}
	})

	// sync jf
	sync_gld_jf(gd.Id, jf_add)

	// match another
	do_match(1000, gd)
}

func sync_gld_jf(gid string, jf_add int32) {
	G.PushAll(gconst.Evt_GsPush_GWarGldJf, []string{gid, core.I32toa(jf_add)}, nil)
}

func send_gld_rewards_mail() {
	L := len(enroll_rank)
	for _, gd := range localdata.Glds {
		// get gld rank
		num := core.ArrayFind(L, func(i int) bool {
			return enroll_rank[i].Id == gd.Id
		})
		if num < 0 {
			continue
		}
		num++

		// rewards
		type rwd_t struct {
			id int32
			n  int64
		}
		var rewards []*rwd_t
		for _, v := range gamedata.ConfGwarRankReward.Items() {
			if v.Rank[0].Low <= int32(num) && int32(num) <= v.Rank[0].High {
				for _, v2 := range v.Reward {
					rewards = append(rewards, &rwd_t{
						id: v2.Id,
						n:  v2.N,
					})
				}
				break
			}
		}

		// send mail
		for plrid, info := range gd.G1Plrs {
			if info.RwdGld {
				continue
			}

			iplr := utils.LoadPlayer(plrid)
			if iplr == nil {
				continue
			}

			m := mail.New(iplr).SetKey(5001)
			m.AddDictInt32("rank", int32(num))

			for _, v := range rewards {
				m.AddAttachment(v.id, float64(v.n))
			}

			m.Send()

			info.RwdGld = true
		}
	}
}

func send_plr_rewards_mail() {
	for _, gd := range localdata.Glds {
		// make plr jf rank
		type plr_rank_t struct {
			id   string
			info *g1_plrinfo_t
		}
		arr := make([]*plr_rank_t, 0, 100)

		for plrid, info := range gd.G1Plrs {
			arr = append(arr, &plr_rank_t{
				id:   plrid,
				info: info,
			})
		}

		sort.Slice(arr, func(i, j int) bool {
			return arr[i].info.Jf > arr[j].info.Jf
		})

		// send mail
		for i, v := range arr {
			if v.info.RwdPlr {
				continue
			}

			iplr := utils.LoadPlayer(v.id)
			if iplr == nil {
				continue
			}

			rk := int32(i) + 1

			m := mail.New(iplr).SetKey(5002)
			m.AddDictInt32("rank", rk)

			conf := gamedata.ConfGwarStaffRankReward.Query(rk)
			if conf != nil {
				for _, v := range conf.Reward {
					m.AddAttachment(v.Id, float64(v.N))
				}
			}

			m.Send()

			v.info.RwdPlr = true
		}
	}
}

// ============================================================================

func (self *localdata_t) reset_data() {
	self.Glds = make(map[string]*gld_data_t)
	self.Enrolled = false
}

// ============================================================================

func (self *gld_data_t) get_g1_plrinfo(id string) *g1_plrinfo_t {
	e := self.G1Plrs[id]
	if e == nil {
		e = &g1_plrinfo_t{}
		self.G1Plrs[id] = e
	}

	return e
}

func (self *gld_data_t) get_g2_plrinfo(id string) *g2_plrinfo_t {
	e := self.G2Plrs[id]
	if e == nil {
		e = &g2_plrinfo_t{}
		self.G2Plrs[id] = e
	}

	return e
}

// ============================================================================

func (self *enroll_base_info_t) tomsg() *msg.GWarGuildBase {
	jf := int32(0)
	e := enroll_index[self.Id]
	if e != nil {
		jf = e.Jf
	}

	return &msg.GWarGuildBase{
		Id:     self.Id,
		Name:   self.Name,
		Icon:   self.Icon,
		Lv:     self.Lv,
		N:      self.N,
		ApMean: self.ApMean,
		ApMax:  self.ApMax,
		SvrId:  self.SvrId,
		Jf:     jf,
	}
}
