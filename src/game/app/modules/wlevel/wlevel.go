package wlevel

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/math"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================
const (
	C_WLevel_FightOneKey = 5
)

// ============================================================================

// 推图
type WLevel struct {
	LvNum int32 // 已通关

	GJTs      time.Time       // 开始挂机时间
	GJLootTs  time.Time       // 上次计算时间
	GJLootRwd map[int32]int64 // 挂机奖励

	plr IPlayer
}

// ============================================================================

func NewWLevel() *WLevel {
	return &WLevel{
		GJTs:      time.Now(),
		GJLootTs:  time.Now(),
		GJLootRwd: make(map[int32]int64),
	}
}

// ============================================================================

func (self *WLevel) Init(plr IPlayer) {
	self.plr = plr
}

func (self *WLevel) team_atkpwr(T *msg.TeamFormation) {
	if !self.plr.IsSetTeam(gconst.TeamType_Dfd) {
		self.plr.SetTeam(gconst.TeamType_Dfd, T)

		self.plr.SendMsg(&msg.GS_SetTeam_R{
			ErrorCode: Err.OK,
			Tp:        gconst.TeamType_Dfd,
			T:         T,
		})
	}

	atk_power := int32(0)
	for seq := range T.Formation {
		hero := self.plr.GetBag().FindHero(seq)
		if hero != nil {
			atk_power += hero.GetAtkPower()
		}
	}

	if float64(atk_power) > self.plr.GetAttainObjVal(gconst.AttainId_MaxAtkPwr) {
		evtmgr.Fire(gconst.Evt_PlrAtkPwr, self.plr, atk_power)
	}
}

func (self *WLevel) Fight(T *msg.TeamFormation, cb func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards)) {
	conf := gamedata.ConfWorldLevelM.Query(self.LvNum + 1)
	if conf == nil {
		cb(Err.Failed, nil, nil)
		return
	}

	if !self.plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil)
		return
	}

	self.team_atkpwr(T)

	// gen input
	T2 := battle.NewMonsterTeam()
	for i, v := range conf.Monster {
		T2.AddMonster(v.Id, v.Lv, int32(i))
	}

	// modify prop
	ratio := conf.PowerRatio
	if self.plr.GetLevel() > conf.PowerSwitch {
		ratio = math.MaxFloat32(0.25,
			conf.PowerRatio*(1-float32(self.plr.GetLevel()-conf.PowerSwitch)*0.05))
	}
	T2.ModifyProps(ratio)

	input := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(T),
		T2: T2.ToMsg_BattleTeam(),
		Args: map[string]string{
			"LvNum":     core.I32toa(self.LvNum + 1),
			"Module":    "WLEVEL",
			"RoundType": core.I32toa(conf.RoundType),
		},
	}

	// fight
	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			cb(Err.Common_BattleResError, nil, nil)
			return
		}

		var rwds *msg.Rewards
		if r.Winner == 1 {
			op := self.plr.GetBag().NewOp(gconst.ObjFrom_WLevelFight)
			for _, v := range conf.Reward {
				op.Inc(v.Id, v.N)
			}

			self.LvNum++

			rwds = op.Apply().ToMsg()

			evtmgr.Fire(gconst.Evt_WLevelLv, self.plr, self.LvNum, int32(1))
		}

		evtmgr.Fire(gconst.Evt_WLevelFight, self.plr, self.LvNum, r.Winner, T)

		replay := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		cb(Err.OK, replay, rwds)
	})
}

// 一键战斗
func (self *WLevel) FightOneKey(T *msg.TeamFormation, cb func(ec int32, replays []*msg.BattleReplay, rwds *msg.Rewards)) {
	if !self.plr.IsModuleOpen(gconst.ModuleId_WLevelFightOneKey) {
		cb(Err.Plr_ModuleLocked, nil, nil)
		return
	}

	if !self.plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil)
		return
	}

	self.team_atkpwr(T)

	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_WLevelOneKeyFight)
	cop.DecCounter(gconst.Cnt_WLevelFightOneKey, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		cb(ec, nil, nil)
		return
	}

	// fight
	var replays []*msg.BattleReplay
	var fight_one func()
	winCnt, fight_idx := int32(0), int32(0)

	fight_one = func() {
		fight_idx++

		curLvNum := self.LvNum + 1 + winCnt
		conf := gamedata.ConfWorldLevelM.Query(curLvNum)
		if conf == nil {
			cb(Err.Failed, nil, nil)
			return
		}

		// gen input
		T2 := battle.NewMonsterTeam()
		for i, v := range conf.Monster {
			T2.AddMonster(v.Id, v.Lv, int32(i))
		}

		// modify prop
		ratio := conf.PowerRatio
		if self.plr.GetLevel() > conf.PowerSwitch {
			ratio = math.MaxFloat32(0.25,
				conf.PowerRatio*(1-float32(self.plr.GetLevel()-conf.PowerSwitch)*0.05))
		}
		T2.ModifyProps(ratio)

		input := &msg.BattleInput{
			T1: self.plr.ToMsg_BattleTeam(T),
			T2: T2.ToMsg_BattleTeam(),
			Args: map[string]string{
				"LvNum":     core.I32toa(curLvNum),
				"Module":    "WLEVEL",
				"RoundType": core.I32toa(conf.RoundType),
			},
		}

		// fight
		battle.Fight(input, func(r *msg.BattleResult) {
			if r == nil {
				cb(Err.Failed, nil, nil)
				return
			}

			if r.Winner == 1 {
				for _, v := range conf.Reward {
					cop.Inc(v.Id, v.N)
				}
				winCnt++
			}

			evtmgr.Fire(gconst.Evt_WLevelLv, self.plr, self.LvNum+winCnt, int32(1))
			evtmgr.Fire(gconst.Evt_WLevelFight, self.plr, self.LvNum+winCnt, r.Winner, T)

			replays = append(replays, &msg.BattleReplay{
				Ts: time.Now().Unix(),
				Bi: input,
				Br: r,
			})

			// 完成,没完成的报错出去
			if fight_idx >= C_WLevel_FightOneKey {
				rwds := cop.Apply().ToMsg()
				self.LvNum += winCnt

				cb(Err.OK, replays, rwds)

				return
			}

			fight_one()
		})
	}

	fight_one()
}

func (self *WLevel) GMSetLevel(v int32) {
	if self.LvNum == v {
		return
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_WLevelFight)
	if self.LvNum < v {
		for i := self.LvNum + 1; i <= v; i++ {
			conf := gamedata.ConfWorldLevelM.Query(i)
			if conf == nil {
				return
			}

			for _, a := range conf.Reward {
				op.Inc(a.Id, a.N)
			}

			self.LvNum = i

			evtmgr.Fire(gconst.Evt_WLevelLv, self.plr, self.LvNum, int32(1))
		}
	} else {
		notEnough := false
		for i := self.LvNum; i >= v+1; i-- {
			conf := gamedata.ConfWorldLevelM.Query(i)
			if conf == nil {
				return
			}

			if !notEnough {
				for _, a := range conf.Reward {
					op.Dec(a.Id, a.N)
				}

				if ec := op.CheckEnough(); ec != Err.OK {
					for _, a := range conf.Reward {
						op.Ret(a.Id, a.N)
					}

					notEnough = true
				}
			}

			self.LvNum = i - 1
		}
	}

	op.Apply().ToMsg()
}

// ============================================================================
// 挂机计算
func (self *WLevel) GJLoot() {
	gjLv := self.LvNum + 1
	maxLv := gamedata.ConfLimitM.Query().MaxWLevelLv
	if gjLv > maxLv {
		gjLv = maxLv
	}

	conf_w := gamedata.ConfWorldLevelM.Query(gjLv)
	if conf_w == nil {
		return
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return
	}

	now := time.Now()

	c_max := self.plr.GetCounter().GetMax(gconst.Cnt_WLevelGJTime)
	max_ts := self.GJTs.Add(time.Hour * time.Duration(int64(conf_g.ExploreTimeLimit)+c_max))
	end_ts := core.MinTime(now, max_ts)

	n := int64(end_ts.Sub(self.GJLootTs).Minutes())
	if n <= 0 {
		return
	}

	for _, v := range conf_w.MinuteCurrency {
		m := float64(v.N) * float64(n) * (1 + float64(self.plr.GetCounter().GetMax(gconst.Cnt_WLevelGJExtReward))/100)
		self.GJLootRwd[v.Id] += int64(m)
	}

	for i := int64(0); i < n; i++ {
		for _, v := range utils.Drop(self.plr, conf_w.ExploreDrop) {
			self.GJLootRwd[v.Id] += v.N
		}
	}

	self.GJLootTs = now.Add(self.GJLootTs.Add(time.Minute * time.Duration(n)).Sub(end_ts))

	// offline full
	if !self.plr.IsOnline() && !max_ts.After(now) {
		evtmgr.Fire(gconst.Evt_OfflineWlevelGjFull, self.plr)
	}
}

// 挂机奖励
func (self *WLevel) GJLootTake() (ec int32, rwds *msg.Rewards) {
	if len(self.GJLootRwd) == 0 {
		return Err.WLevel_GJEmpty, nil
	}

	self.GJLoot()

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_WLevelGJ)

	for id, n := range self.GJLootRwd {
		op.Inc(id, n)
	}

	now := time.Now()

	self.GJTs = now
	self.GJLootTs = now
	self.GJLootRwd = make(map[int32]int64)

	rwds = op.Apply().ToMsg()

	evtmgr.Fire(gconst.Evt_WLevelGj, self.plr)

	return Err.OK, rwds
}

// ============================================================================

func (self WLevel) ToMsg() *msg.WLevelData {
	return &msg.WLevelData{
		LvNum:     self.LvNum,
		GJTs:      self.GJTs.Unix(),
		GJLootRwd: self.GJLootRwd,
	}
}
