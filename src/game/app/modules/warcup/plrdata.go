package warcup

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 本服杯赛-任务
type WarCup struct {
	VerTs time.Time // 当前版本

	Attain    map[int32]*attain_obj_t // [id]item
	TaskTaken []int32                 // 已经领取任务

	plr IPlayer
}

type attain_obj_t struct {
	Id  int32
	Val float64 // progress value
}

// ============================================================================

func NewWarCup() *WarCup {
	return &WarCup{
		Attain: make(map[int32]*attain_obj_t),
	}
}

// ============================================================================

func (self *WarCup) Init(plr IPlayer) {
	self.plr = plr
}

func (self *WarCup) get(oid int32) *attain_obj_t {
	obj := self.Attain[oid]
	if obj == nil {
		self.Attain[oid] = &attain_obj_t{Id: oid}
		return self.Attain[oid]
	}

	return obj
}

func (self *WarCup) CheckVersion() {
	if !self.VerTs.IsZero() && self.VerTs.Equal(g_t0) {
		return
	}

	//reset
	self.VerTs = g_t0
	self.Attain = make(map[int32]*attain_obj_t)
	self.TaskTaken = []int32{}
}

func (self *WarCup) Take(id int32) (int32, *msg.Rewards) {
	for _, v := range self.TaskTaken {
		if v == id {
			return Err.Plr_TakenBefore, nil
		}
	}

	conf := gamedata.ConfWarCupGuessTask.Query(id)
	if conf == nil || len(conf.GuessTaskAttain) == 0 {
		return Err.Failed, nil
	}

	for _, v := range conf.GuessTaskAttain {
		val := self.get(v.AttainId).Val

		if val == 0 || val < v.P2 {
			return Err.Plr_CondLimited, nil
		}
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_WarCup)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.TaskTaken = append(self.TaskTaken, id)

	return Err.OK, rwds
}

// ============================================================================

func (self *WarCup) ToMsg() *msg.WarCupData {
	self.CheckVersion()

	ret := WarCup_ToMsg(self.plr)

	for _, v := range self.Attain {
		ret.Attain = append(ret.Attain, &msg.WarCupAttainObj{
			OId: v.Id,
			Val: v.Val,
		})
	}

	ret.TaskTaken = self.TaskTaken

	return ret
}

// ============================================================================
// implements ICondObj interface

func (self *attain_obj_t) GetVal() float64 {
	return self.Val
}

func (self *attain_obj_t) SetVal(v float64) {
	self.Val = v
}

func (self *attain_obj_t) AddVal(v float64) {
	self.Val += v
}

func (self *attain_obj_t) Done(body interface{}, confid int32, isChange bool) {
	if !isChange {
		return
	}

	body.(IPlayer).SendMsg(&msg.GS_WarCupAttainObjValueChanged{
		OId: self.Id,
		Val: self.Val,
	})
}
