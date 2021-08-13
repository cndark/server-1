package taskmonth

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 每月任务
type TaskMonth struct {
	Items item_m

	plr IPlayer
}

type item_m map[int32]*item_t
type item_t struct {
	Val   float64
	Fin   bool
	Taken bool
}

// ============================================================================

func NewTaskMonth() *TaskMonth {
	return &TaskMonth{
		Items: make(item_m),
	}
}

// ============================================================================

func (self *TaskMonth) Init(plr IPlayer) {
	self.plr = plr
}

func (self *TaskMonth) IsCompleted(id int32) bool {
	item := self.Items[id]
	return item != nil && item.Fin

}

func (self *TaskMonth) get(id int32) (item *item_t) {
	item = self.Items[id]
	if item == nil {
		item = &item_t{}
		self.Items[id] = item
	}
	return
}

func (self *TaskMonth) reset_month() {
	for id, _ := range self.Items {
		self.Items[id].Val = 0
		self.Items[id].Fin = false
		self.Items[id].Taken = false
	}
}

func (self *TaskMonth) Take(id int32) (int32, *msg.Rewards) {
	conf := gamedata.ConfTaskMonth.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	item := self.Items[id]

	if item == nil {
		return Err.TaskMonth_NotFound, nil
	}

	if item.Taken {
		return Err.TaskMonth_AlreadRewarded, nil
	}

	if !item.Fin {
		return Err.TaskMonth_NotCompleted, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TaskMonth)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	item.Taken = true

	rwd := op.Apply().ToMsg()

	return Err.OK, rwd
}

func (self *TaskMonth) ToMsg() *msg.TaskMonthData {

	ret := &msg.TaskMonthData{}

	for id, item := range self.Items {
		if gamedata.ConfTaskMonth.Query(id) == nil {
			continue
		}
		ret.Items = append(ret.Items, &msg.TaskMonthItem{
			Id:  id,
			Val: item.Val,
			Fin: item.Fin,
			T:   item.Taken,
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
	conf := gamedata.ConfTaskMonth.Query(confid)
	if conf == nil {
		return
	}

	dtask := body.(IPlayer).GetTaskMonth()
	if dtask.IsCompleted(confid) {
		return
	}

	if self.Val < conf.P2 {
		body.(IPlayer).SendMsg(&msg.GS_TaskMonthValueChanged{
			Id:  confid,
			Val: self.Val,
		})
	} else {
		self.Fin = true

		body.(IPlayer).SendMsg(&msg.GS_TaskMonthItemCompleted{
			Id: confid,
		})

		// evt
		evtmgr.Fire(gconst.Evt_TaskMonth_Fin, body.(IPlayer), confid, conf.Type)
	}
}
