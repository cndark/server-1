package tower

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/math"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

// 爬塔
type Tower struct {
	LvNum  int32     // 当前通过层数
	LastTs time.Time // 上次领取时间

	plr IPlayer
}

// ============================================================================

func NewTower() *Tower {
	return &Tower{
		LastTs: time.Unix(0, 0),
	}
}

// ============================================================================

func (self *Tower) Init(plr IPlayer) {
	self.plr = plr
}

// 战斗
func (self *Tower) Fight(T *msg.TeamFormation, cb func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards)) {
	if !self.plr.IsModuleOpen(gconst.ModuleId_Tower) {
		cb(Err.Plr_ModuleLocked, nil, nil)
		return
	}

	conf := gamedata.ConfTowerM.Query(self.LvNum + 1)
	if conf == nil {
		cb(Err.Failed, nil, nil)
		return
	}

	if !self.plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil)
		return
	}

	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_TowerFight)
	cop.DecCounter(gconst.Cnt_TowerFight, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		cb(ec, nil, nil)
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
		ratio = math.MaxFloat32(0.5,
			conf.PowerRatio*(1-float32(self.plr.GetLevel()-conf.PowerSwitch)*0.1))
	}
	T2.ModifyProps(ratio)

	input := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(T),
		T2: T2.ToMsg_BattleTeam(),
		Args: map[string]string{
			"LvNum":     core.I32toa(self.LvNum + 1),
			"Module":    "TOWER",
			"RoundType": core.I32toa(conf.RoundType),
		},
	}

	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			cb(Err.Common_BattleResError, nil, nil)
			return
		}

		cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_TowerFight)
		var rwds *msg.Rewards
		if r.Winner == 1 {
			for _, v := range conf.Reward {
				cop.Inc(v.Id, v.N)
			}

			self.LvNum++

			evtmgr.Fire(gconst.Evt_TowerLv, self.plr, self.LvNum)
		} else {
			cop.DecCounter(gconst.Cnt_TowerFight, 1)
			if ec := cop.CheckEnough(); ec != Err.OK {
				cb(ec, nil, nil)
				return
			}
		}

		rwds = cop.Apply().ToMsg()

		replay := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		TowerMgr.add_replay(self.LvNum, replay)

		cb(Err.OK, replay, rwds)

		evtmgr.Fire(gconst.Evt_TowerFight, self.plr, self.LvNum, r.Winner, T)
	})
}

// 扫荡
func (self *Tower) Raid() (int32, *msg.Rewards) {
	conf := gamedata.ConfTowerM.Query(self.LvNum)
	if conf == nil {
		return Err.Failed, nil
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return Err.Failed, nil
	}

	now := time.Now()
	if core.IsSameDay(self.LastTs, now) {
		return Err.Plr_TakenBefore, nil
	}

	// rwds
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TowerRaid)

	for _, v := range conf.RaidReward {
		if rand.Float32() < v.Odd {
			op.Inc(v.Id, v.N)
		}
	}

	rwds := op.Apply().ToMsg()

	self.LastTs = now

	evtmgr.Fire(gconst.Evt_TowerRaid, self.plr, self.LvNum)

	return Err.OK, rwds
}

func (self *Tower) GMSetLevel(v int32) {
	if self.LvNum == v {
		return
	}

	op := self.plr.GetCounter().NewOp(gconst.ObjFrom_TowerFight)
	if self.LvNum < v {
		for i := self.LvNum + 1; i <= v; i++ {
			conf := gamedata.ConfTowerM.Query(i)
			if conf == nil {
				return
			}

			for _, a := range conf.Reward {
				op.Inc(a.Id, a.N)
			}

			self.LvNum = i

			evtmgr.Fire(gconst.Evt_TowerLv, self.plr, self.LvNum)
		}
	} else {
		notEnough := false
		for i := self.LvNum; i >= v+1; i-- {
			conf := gamedata.ConfTowerM.Query(i)
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

func (self *Tower) ToMsg() *msg.TowerData {
	return &msg.TowerData{
		LvNum:  self.LvNum,
		LastTs: self.LastTs.Unix(),
	}
}
