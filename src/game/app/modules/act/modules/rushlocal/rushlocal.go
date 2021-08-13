package rushlocal

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/guild"
	"fw/src/game/app/modules/ranksvc"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

var actRushLocal = map[string]*act_t{
	gconst.ActName_RushLocal1: new(act_t),
	gconst.ActName_RushLocal2: new(act_t),
	gconst.ActName_RushLocal3: new(act_t),
	gconst.ActName_RushLocal4: new(act_t),
}

// ============================================================================

func init() {
	for actName := range actRushLocal {
		act.RegisterAct(actName, actRushLocal[actName])
	}

	on_evt()
}

// ============================================================================

type act_t struct {
	act.BaseAct

	rank_tid *core.Timer
}

type data_plr_t struct {
	Take map[int32]bool
}

type data_svr_t struct {
	Ranks rank_m
	valid map[int32]bool
}

type rank_m map[int32]*rank_one_t
type rank_one_t struct {
	Data        map[string]float64 // [plrid 或 guildid] score
	CachedScore map[string]float64 // 活动前最高排行值备份
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return &data_svr_t{
		Ranks: make(rank_m),
		valid: make(map[int32]bool),
	}
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{
		Take: make(map[int32]bool),
	}
}

func (self *act_t) GetSvrData() *data_svr_t {
	return self.GetActRawData().(*data_svr_t)
}

func (self *act_t) GetPlrData(plr IPlayer) *data_plr_t {
	return plr.GetActRawData(self.GetName()).(*data_plr_t)
}

// ============================================================================

func (self *act_t) Started() bool {
	return self.GetStage() == "start"
}

func (self *act_t) Ended() bool {
	return self.GetStage() == "end"
}

func (self *act_t) Closed() bool {
	return self.GetStage() == "close"
}

// ============================================================================

func (self *act_t) OnInit() {
}

func (self *act_t) OnStage() {
	if !self.Closed() {
		self.UpdateRankIdValid()
	}

	if self.Started() {
		self.pick_cached_score()

		self.stop_timer_push_rank()
		self.push_rank()
		self.start_timer_push_rank()
	}

	if self.Ended() {
		self.stop_timer_push_rank()
		// 最后排一次
		loop.SetTimeout(time.Now().Add(time.Duration(1)*time.Minute), func() {
			self.push_rank()
		})
	}

	if self.Closed() {
		svr_data := self.GetSvrData()
		if svr_data != nil {
			svr_data.valid = make(map[int32]bool)
			svr_data.Ranks = make(rank_m)
		}
	}
}

// ============================================================================

func (self *act_t) start_timer_push_rank() {
	self.rank_tid = loop.SetTimeout(time.Now().Add(time.Duration(5)*time.Minute), func() {
		if self.Started() {
			self.push_rank()
		}

		self.start_timer_push_rank()
	})
}

func (self *act_t) stop_timer_push_rank() {
	if self.rank_tid != nil {
		loop.CancelTimer(self.rank_tid)
		self.rank_tid = nil
	}
}

func (self *act_t) push_rank() {
	svr_data := self.GetSvrData()
	if svr_data == nil {
		return
	}

	for rk := range svr_data.valid {
		raw := &ranksvc.RankRaw{
			Id:        rk,
			SortLevel: ranksvc.SortLevel_Local,
		}

		A := []*ranksvc.RankRow{}
		r := svr_data.get_rank(rk)
		if r != nil {
			A = make([]*ranksvc.RankRow, 0, len(r.Data))
			for uid, v := range r.Data {
				if gconst.IsGuildRankType(rk) {
					gld := guild.GuildMgr.FindGuild(uid)
					if gld != nil {
						A = append(A, &ranksvc.RankRow{
							Score: v,
							Info: &ranksvc.RankRowInfo{
								Gld: gld.ToMsg_Row(),
							},
						})
					}
				} else {
					plr := find_player(uid)
					if plr != nil && !plr.IsBan() {
						A = append(A, &ranksvc.RankRow{
							Score: v,
							Info: &ranksvc.RankRowInfo{
								Plr: plr.ToMsg_SimpleInfo(),
							},
						})
					}
				}
			}
		}

		raw.A = A
		raw.ScoreFunc = func(i int) float64 { return A[i].Score }
		raw.InfoFunc = func(i int) *ranksvc.RankRowInfo { return A[i].Info }

		ranksvc.Push(raw)
	}
}

// ============================================================================

// 更新当前存在的榜id
func (self *act_t) UpdateRankIdValid() {
	svr_data := self.GetSvrData()
	if svr_data == nil {
		return
	}

	items := gamedata.ConfActRushLocalM.QueryItems(self.GetConfGrp())
	for _, v := range items {
		svr_data.valid[v.RankId] = true

		svr_data.get_rank(v.RankId)
	}
}

func (self *act_t) IsRankValid(rk int32) bool {
	svr_data := self.GetSvrData()
	if svr_data != nil {
		return svr_data.valid[rk]
	}

	return false
}

// 活动开始时, 缓存历史最大值
func (self *act_t) pick_cached_score() {
	// svr_data := self.GetSvrData()
	// if svr_data == nil {
	// 	return
	// }
}

// 领取排行奖励
func (self *act_t) TakeRewards(plr IPlayer, rk int32, f func(ec int32, rwds *msg.Rewards)) {
	if !self.IsRankValid(rk) {
		f(Err.Activity_RankNotFound, nil)
		return
	}

	if !self.Ended() {
		f(Err.Act_StageError, nil)
		return
	}

	if time.Now().Before(self.GetT1().Add(time.Duration(5) * time.Minute)) {
		f(Err.Activity_TimeLimit, nil)
		return
	}

	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		f(Err.Act_ActPlrDataNotFound, nil)
		return
	}

	if plr_data.Take[rk] {
		f(Err.Activity_TakeBefore, nil)
		return
	}

	svr_data := self.GetSvrData()
	if svr_data == nil {
		f(Err.Act_ActSvrDataNotFound, nil)
		return
	}

	r := svr_data.Ranks[rk]
	if r == nil {
		f(Err.Activity_RankNotFound, nil)
		return
	}

	uid := plr.GetId()
	if gconst.IsGuildRankType(rk) {
		uid = plr.GetGuildId()
		// if uid == "" {
		// 	f(Err.Guild_NotFound, nil)
		// 	return
		// }
	}

	// 没参加不给奖励
	// if r.Data[uid] <= 0 {
	// 	f(Err.Activity_NotJoin, nil)
	// 	return
	// }

	ranksvc.Get(ranksvc.RankType_Local, config.CurGame.Id, rk, func(rows []*ranksvc.RankRow) {
		idx := ranksvc.GetRowPos(rows, uid)
		op := plr.GetBag().NewOp(gconst.ObjFrom_ActRushLocal)

		items := gamedata.ConfActRushLocalM.QueryItems(self.GetConfGrp())
		for _, conf_r := range items {
			if conf_r.RankId == rk {
				isAward := false
				zero_item := make(map[int32]int64)
				for _, conf := range gamedata.ConfActRushLocalRewardM.QueryItems(conf_r.RewardGrp) {
					if len(conf.Rank) > 0 && conf.Rank[0].Low == 0 {
						for _, v := range conf.Reward {
							zero_item[v.Id] = v.N
						}
					}

					if len(conf.Rank) > 0 && idx >= conf.Rank[0].Low && idx <= conf.Rank[0].High {
						for _, v := range conf.Reward {
							op.Inc(v.Id, v.N)
						}
						isAward = true
						goto LOOP
					}
				}

				if !isAward {
					for id, n := range zero_item {
						op.Inc(id, n)
					}
				}
				break
			}
		}

	LOOP:
		plr_data.Take[rk] = true

		rwds := op.Apply().ToMsg()
		f(Err.OK, rwds)
	})

	return
}

func (self *act_t) ToMsg(plr IPlayer, rk int32) (ret []*msg.ActRushLocalRankData) {
	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return
	}

	svr_data := self.GetSvrData()
	if svr_data == nil {
		return
	}

	var rks []int32
	if rk != 0 {
		if !self.IsRankValid(rk) {
			return
		}
		rks = append(rks, rk)
	} else {
		for k := range svr_data.valid {
			rks = append(rks, k)
		}
	}

	for _, k := range rks {
		uid := plr.GetId()
		if gconst.IsGuildRankType(k) {
			uid = plr.GetGuildId()
		}

		take := plr_data.Take[k]

		r := svr_data.get_rank(k)
		score := r.Data[uid]

		ret = append(ret, &msg.ActRushLocalRankData{
			ActName:     self.GetName(),
			RankId:      k,
			Take:        take,
			SelfScore:   score,
			CachedScore: r.CachedScore[uid],
		})
	}

	return
}

// ============================================================================

func (self *data_svr_t) get_rank(rk int32) *rank_one_t {
	r := self.Ranks[rk]
	if r == nil {
		r = &rank_one_t{
			Data:        make(map[string]float64),
			CachedScore: make(map[string]float64),
		}

		self.Ranks[rk] = r
	} else if r.Data == nil {
		r.Data = make(map[string]float64)
	} else if r.CachedScore == nil {
		r.CachedScore = make(map[string]float64)
	}

	return r
}

func (self *data_svr_t) update_rank_rec(isAdd bool, rk int32, uid string, v float64) {
	r := self.get_rank(rk)
	if isAdd {
		r.Data[uid] += v
	} else { // cached
		add := v - r.CachedScore[uid]
		if r.Data[uid] < add {
			r.Data[uid] = add
		}
	}
}

// ============================================================================
// api

func GetInfo(plr IPlayer, actName string, rk int32) (ret []*msg.ActRushLocalRankData) {
	rush, ok := actRushLocal[actName]
	if !ok {
		for _, rush := range actRushLocal {
			if !rush.Closed() {
				ret = append(ret, rush.ToMsg(plr, rk)...)
			}
		}
	} else {
		ret = append(ret, rush.ToMsg(plr, rk)...)
	}

	return
}

func TakeRankRewards(plr IPlayer, actName string, rk int32) {
	f := func(ec int32, rwds *msg.Rewards) {
		res := &msg.GS_ActRushLocalTake_R{}

		res.ErrorCode = ec
		res.ActName = actName
		res.RankId = rk
		res.Rewards = rwds

		plr.SendMsg(res)
	}

	rush, ok := actRushLocal[actName]
	if ok {
		rush.TakeRewards(plr, rk, f)
	} else {
		f(Err.Act_ActNotFound, nil)
	}
}

// ============================================================================

func on_evt() {
	// load gamedata
	evtmgr.On(gconst.Evt_GameDataReload, func(args ...interface{}) {
		names := args[0].([]string)
		b := false
		for _, name := range names {
			if name == "actRushLocal" {
				b = true
				break
			}
		}

		if b {
			for _, rush := range actRushLocal {
				if rush.Closed() {
					continue
				}

				rush.UpdateRankIdValid()
			}
		}
	})

	// wlevel
	evtmgr.On(gconst.Evt_WLevelLv, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		lvNum := args[1].(int32)
		addNum := args[2].(int32)

		if addNum > 0 {
			for _, rush := range actRushLocal {
				if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_WLevel) {
					continue
				}

				svr_data := rush.GetSvrData()
				if svr_data != nil {
					svr_data.update_rank_rec(false, gconst.RankId_RushLocal_WLevel, plr.GetId(), float64(lvNum))
				}
			}
		}
	})

	// guild boss
	evtmgr.On(gconst.Evt_GuildBossFight, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		dmg := args[2].(float64)

		for _, rush := range actRushLocal {
			if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_GuildBossDmg) {
				continue
			}

			svr_data := rush.GetSvrData()
			if svr_data != nil {
				svr_data.update_rank_rec(true, gconst.RankId_RushLocal_GuildBossDmg, plr.GetId(), float64(dmg))
			}
		}
	})

	// atkpower
	evtmgr.On(gconst.Evt_PlrAtkPwr, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		atkpwr := args[1].(int32)

		for _, rush := range actRushLocal {
			if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_AtkPwr) {
				continue
			}

			svr_data := rush.GetSvrData()
			if svr_data != nil {
				svr_data.update_rank_rec(false, gconst.RankId_RushLocal_AtkPwr, plr.GetId(), float64(atkpwr))
			}
		}
	})

	// tower
	evtmgr.On(gconst.Evt_TowerLv, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		lvNum := args[1].(int32)

		for _, rush := range actRushLocal {
			if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_Tower) {
				continue
			}

			svr_data := rush.GetSvrData()
			if svr_data != nil {
				svr_data.update_rank_rec(false, gconst.RankId_RushLocal_Tower, plr.GetId(), float64(lvNum))
			}
		}
	})

	// draw
	evtmgr.On(gconst.Evt_Draw, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		tp := args[1].(string)
		n := args[3].(int32)

		conf_g := gamedata.ConfActivityPublic.Query(1)
		if conf_g == nil {
			return
		}

		score := int32(0)
		for _, v := range conf_g.RushLocalDrawScore {
			if v.Tp == tp {
				score = v.Sc
				break
			}
		}

		if score == 0 {
			return
		}

		for _, rush := range actRushLocal {
			if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_Draw) {
				continue
			}

			svr_data := rush.GetSvrData()
			if svr_data != nil {
				svr_data.update_rank_rec(true, gconst.RankId_RushLocal_Draw, plr.GetId(), float64(n*score))
			}
		}
	})

	// marvelroll
	evtmgr.On(gconst.Evt_MarvelRoll, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		grp := args[1].(string)
		n := args[3].(int32)

		conf_g := gamedata.ConfActivityPublic.Query(1)
		if conf_g == nil {
			return
		}

		score := int32(0)
		for _, v := range conf_g.RushLocalMarvelRollScore {
			if v.Tp == grp {
				score = v.Sc
				break
			}
		}

		if score == 0 {
			return
		}

		for _, rush := range actRushLocal {
			if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_MarvelRoll) {
				continue
			}

			svr_data := rush.GetSvrData()
			if svr_data != nil {
				svr_data.update_rank_rec(true, gconst.RankId_RushLocal_MarvelRoll, plr.GetId(), float64(n*score))
			}
		}
	})

	// arena
	evtmgr.On(gconst.Evt_ArenaFight, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		isWin := args[1].(bool)

		conf := gamedata.ConfActivityPublic.Query(1)
		if conf == nil || len(conf.RushLocalArenaScore) < 2 {
			return
		}

		score := conf.RushLocalArenaScore[0]
		if isWin {
			score = conf.RushLocalArenaScore[1]
		}

		for _, rush := range actRushLocal {
			if !rush.Started() || !rush.IsRankValid(gconst.RankId_RushLocal_Arena) {
				continue
			}

			svr_data := rush.GetSvrData()
			if svr_data != nil {
				svr_data.update_rank_rec(true, gconst.RankId_RushLocal_Arena, plr.GetId(), float64(score))
			}
		}
	})

}
