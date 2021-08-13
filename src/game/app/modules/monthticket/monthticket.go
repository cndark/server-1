package monthticket

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
)

// ============================================================================
// 月票
type MonthTicket struct {
	Items  item_m
	PickId []int32 // 今天的任务id

	IsBuy      bool    // 是否购买了月票
	Lv         int32   // 等级
	Exp        int32   // 经验
	TakeBase   int32   // 领取基础奖励位置
	TakeTicket int32   // 领取月票奖励位置
	TaskTaken  []int32 // 已经领取的任务

	plr IPlayer
}

type item_m map[int32]*item_t
type item_t struct {
	Val float64
	Fin bool
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		pid := args[1].(int32)

		conf_prod := gamedata.ConfBillProduct.Query(pid)
		if conf_prod == nil || conf_prod.TypeId != gconst.Bill_MonthTicket {
			return
		}

		if conf_prod.PayId != gconst.Bill_PayId_MonthTicket {
			return
		}

		plr.GetMonthTicket().IsBuy = true

		plr.SendMsg(&msg.GS_MonthTicketIsBuy{
			IsBuy: true,
		})
	})

	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr.GetMonthTicket().daily_reset()
	})

	evtmgr.On(gconst.Evt_PlrResetMonthly, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr.GetMonthTicket().month_reset()
	})

}

func NewMonthTicket() *MonthTicket {
	return &MonthTicket{
		Items:      make(item_m),
		TakeBase:   int32(-1),
		TakeTicket: int32(-1),
	}
}

// ============================================================================

func (self *MonthTicket) Init(plr IPlayer) {
	self.plr = plr
}

func (self *MonthTicket) get(id int32) (item *item_t) {
	item = self.Items[id]
	if item == nil {
		item = &item_t{}
		self.Items[id] = item
	}

	return
}

func (self *MonthTicket) daily_reset() {
	for _, item := range self.Items {
		item.Val = 0
		item.Fin = false
	}

	self.TaskTaken = []int32{}

	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf != nil && conf.MonthTicketTaskNum > 0 {
		slt := []int32{}
		for _, v := range gamedata.ConfMonthTicketTask.Items() {
			slt = append(slt, v.Id)
		}

		L := len(gamedata.ConfMonthTicketTask.Items())
		idx := rand.Perm(int(L))

		n := int32(0)
		self.PickId = []int32{}
		for _, i := range idx {
			if n >= conf.MonthTicketTaskNum {
				break
			}

			self.PickId = append(self.PickId, slt[i])
			n++
		}
	}
}

func (self *MonthTicket) month_reset() {
	self.IsBuy = false
	self.Lv = 0
	self.Exp = 0
	self.TakeBase = int32(-1)
	self.TakeTicket = int32(-1)
}

func (self *MonthTicket) isCompleted(id int32) bool {
	item := self.Items[id]
	return item != nil && item.Fin
}

// 一键领取赏金
func (self *MonthTicket) TakeOneKey() (int32, *msg.Rewards) {
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_MonthTicket)

	for lv := self.TakeBase + 1; lv <= self.Lv; lv++ {
		conf := gamedata.ConfMonthTicket.Query(lv)
		if conf != nil {
			for _, v := range conf.BaseReward {
				op.Inc(v.Id, v.N)
			}
		}
	}

	self.TakeBase = self.Lv

	if self.IsBuy {
		for lv := self.TakeTicket + 1; lv <= self.Lv; lv++ {
			conf := gamedata.ConfMonthTicket.Query(lv)
			if conf != nil {
				for _, v := range conf.TicketReward {
					op.Inc(v.Id, v.N)
				}
			}
		}

		self.TakeTicket = self.Lv
	}

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// 领取赏金
func (self *MonthTicket) Take(lv int32) (int32, *msg.Rewards) {
	if lv > self.Lv {
		return Err.Failed, nil
	}

	if lv <= self.TakeBase {
		if !self.IsBuy {
			return Err.Plr_TakenBefore, nil
		} else if lv <= self.TakeTicket {
			return Err.Plr_TakenBefore, nil
		}
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_MonthTicket)

	if lv > self.TakeBase {
		conf := gamedata.ConfMonthTicket.Query(self.TakeBase + 1)
		if conf != nil {
			for _, v := range conf.BaseReward {
				op.Inc(v.Id, v.N)
			}
			self.TakeBase++
		}
	}

	if self.IsBuy && lv > self.TakeTicket {
		conf := gamedata.ConfMonthTicket.Query(self.TakeTicket + 1)
		if conf != nil {
			for _, v := range conf.TicketReward {
				op.Inc(v.Id, v.N)
			}
			self.TakeTicket++
		}
	}

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// 购买直升
func (self *MonthTicket) BuyUp(n int32) int32 {
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_MonthTicket)

	cnt := int32(0)
	for i := int32(0); i < n; i++ {
		conf := gamedata.ConfMonthTicket.Query(self.Lv + i)
		if conf == nil || conf.Exp == 0 {
			break
		}

		cnt++
		op.Dec(gconst.Diamond, conf.UpCost)
	}

	if cnt == 0 {
		return Err.Failed
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec
	}

	self.Lv += cnt
	self.Exp = 0

	op.Apply()

	return Err.OK
}

// 任务领取
func (self *MonthTicket) TaskTake(id int32) int32 {
	conf := gamedata.ConfMonthTicketTask.Query(id)
	if conf == nil {
		return Err.Failed
	}

	for _, v := range self.TaskTaken {
		if v == id {
			return Err.Plr_TakenBefore
		}
	}

	// auto take
	isFound := false
	for _, v := range self.PickId {
		if v == id {
			isFound = true
			break
		}
	}

	if !isFound {
		return Err.MonthTicket_NotFound
	}

	conf_l := gamedata.ConfMonthTicket.Query(self.Lv)
	if conf_l == nil || conf_l.Exp <= 0 {
		return Err.Failed
	}

	self.Exp += conf.GetExp
	if self.Exp >= conf_l.Exp {
		self.Exp -= conf_l.Exp
		self.Lv++
	}

	self.TaskTaken = append(self.TaskTaken, id)

	return Err.OK
}

// ============================================================================

func (self *MonthTicket) ToMsg() *msg.MonthTicketData {
	ret := &msg.MonthTicketData{
		IsBuy:      self.IsBuy,
		Lv:         self.Lv,
		Exp:        self.Exp,
		TakeBase:   self.TakeBase,
		TakeTicket: self.TakeTicket,
		PickId:     self.PickId,
		TaskTaken:  self.TaskTaken,
	}

	for id, v := range self.Items {
		ret.Items = append(ret.Items, &msg.MonthTicketItem{
			Id:  id,
			Fin: v.Fin,
			Val: v.Val,
		})
	}

	return ret
}

// ============================================================================
// implements ICondVal interface

func (self *item_t) GetVal() float64 {
	return self.Val
}

func (self *item_t) SetVal(v float64) {
	self.Val = v
}

func (self *item_t) AddVal(v float64) {
	self.Val += v
}

func (self *item_t) Done(body interface{}, confid int32, isChange bool) {
	conf := gamedata.ConfMonthTicketTask.Query(confid)
	if conf == nil {
		return
	}

	task := body.(IPlayer).GetMonthTicket()
	if task.isCompleted(confid) {
		return
	}

	if self.Val < conf.P2 {
		body.(IPlayer).SendMsg(&msg.GS_MonthTicketValueChanged{
			Id:  confid,
			Val: self.Val,
		})
	} else {
		self.Fin = true

		body.(IPlayer).SendMsg(&msg.GS_MonthTicketItemCompleted{
			Id: confid,
		})
	}
}
