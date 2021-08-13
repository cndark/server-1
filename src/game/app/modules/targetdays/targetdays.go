package targetdays

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 开服庆典(N日目标)
type TargetDays struct {
	Taken  []int32         // 已经领取
	BuyCnt map[int32]int32 // 购买次数

	plr IPlayer
}

// ============================================================================

func NewTargetDays() *TargetDays {
	return &TargetDays{
		BuyCnt: make(map[int32]int32),
	}
}

// ============================================================================

func (self *TargetDays) Init(plr IPlayer) {
	self.plr = plr
}

func (self *TargetDays) IsComplete(id int32) bool {
	for _, v := range self.Taken {
		if v == id {
			return true
		}
	}

	return false
}

// 领取任务
func (self *TargetDays) Take(id int32) (int32, *msg.Rewards) {
	if self.IsComplete(id) {
		return Err.Plr_TakenBefore, nil
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return Err.Failed, nil
	}

	conf_t := gamedata.ConfTargetDays.Query(id)
	if conf_t == nil || len(conf_t.AttainTab) == 0 {
		return Err.Failed, nil
	}

	curDay := core.DistanceDays(self.plr.GetCreateTs())
	if curDay > conf_g.TargetDaysDur {
		return Err.TargetDays_Closed, nil
	}

	if curDay < conf_t.Day {
		return Err.Common_TimeNotUp, nil
	}

	for _, v := range conf_t.AttainTab {
		val := self.plr.GetAttainObjVal(v.AttainId)

		if val == 0 || val < v.P2 {
			return Err.TargetDays_NotCompleted, nil
		}
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TargetDays)

	for _, v := range conf_t.Reward {
		op.Inc(v.Id, v.N)
	}

	self.Taken = append(self.Taken, id)

	rwds := op.Apply().ToMsg()

	// evt
	evtmgr.Fire(gconst.Evt_TargetDays_Take, self.plr, id, conf_t.Type)

	return Err.OK, rwds
}

// 折扣购买
func (self *TargetDays) Buy(id int32) (int32, *msg.Rewards) {
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return Err.Failed, nil
	}

	curDay := core.DistanceDays(self.plr.GetCreateTs())
	if curDay > conf_g.TargetDaysDur {
		return Err.TargetDays_Closed, nil
	}

	conf_t := gamedata.ConfTargetDaysBuy.Query(id)
	if conf_t == nil {
		return Err.Failed, nil
	}

	cnt := self.BuyCnt[id]
	if cnt >= conf_t.MaxBuy {
		return Err.TargetDays_BuyCntMax, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TargetDays)
	op.Dec(gconst.Diamond, conf_t.Price)
	if ec := op.CheckEnough(); ec != Err.OK {
		return ec, nil
	}

	for _, v := range conf_t.Item {
		op.Inc(v.Id, v.N)
	}

	self.BuyCnt[id]++

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// ============================================================================

func (self *TargetDays) ToMsg() *msg.TargetDaysData {
	ret := &msg.TargetDaysData{
		Taken:  self.Taken,
		BuyCnt: self.BuyCnt,
	}

	return ret
}
