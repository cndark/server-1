package taskdaily

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 日常任务
type TaskDaily struct {
	Items   item_m
	BoxTake []int32

	plr IPlayer
}

type item_m map[int32]*item_t
type item_t struct {
	Val   float64
	Fin   bool
	Taken bool
}

// ============================================================================
func init() {
	// reset daily
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		plr.GetTaskDaily().reset_daily()
	})
}

// ============================================================================

func NewTaskDaily() *TaskDaily {
	return &TaskDaily{
		Items: make(item_m),
	}
}

func (self *TaskDaily) Init(plr IPlayer) {
	self.plr = plr
}

func (self *TaskDaily) reset_daily() {
	for _, item := range self.Items {
		item.Val = 0
		item.Fin = false
		item.Taken = false
	}

	self.BoxTake = self.BoxTake[:0]
}

func (self *TaskDaily) get(id int32) (item *item_t) {
	item = self.Items[id]
	if item == nil {
		item = &item_t{}
		self.Items[id] = item
	}

	return
}

func (self *TaskDaily) IsCompleted(id int32) bool {
	item := self.Items[id]
	return item != nil && item.Fin
}

func (self *TaskDaily) TakeBoxReward(id int32) (ec int32, rwd *msg.Rewards) {
	for _, v := range self.BoxTake {
		if v == id {
			ec = Err.TaskDaily_AlreadRewarded
			return
		}
	}

	conf := gamedata.ConfTaskDailyBox.Query(id)
	if conf == nil {
		ec = Err.Failed
		return
	}

	fin := int32(0)
	for _, v := range self.Items {
		if v.Fin {
			fin++
		}
	}

	if fin < conf.ActiveCond {
		ec = Err.TaskDaily_ActiveNotEnough
		return
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TaskDaily)
	for _, v := range conf.Reward {
		op.Inc(v.Id, int64(v.N))
	}

	self.BoxTake = append(self.BoxTake, id)

	rwd = op.Apply().ToMsg()

	ec = Err.OK
	return
}

func (self *TaskDaily) ToMsg() *msg.TaskDailyData {
	ret := &msg.TaskDailyData{
		BoxTake: self.BoxTake,
	}

	for id, item := range self.Items {
		//if table unexist
		if gamedata.ConfTaskDaily.Query(id) == nil {
			continue
		}

		ret.Items = append(ret.Items, &msg.TaskDailyItem{
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
	conf := gamedata.ConfTaskDaily.Query(confid)
	if conf == nil {
		return
	}

	dtask := body.(IPlayer).GetTaskDaily()
	if dtask.IsCompleted(confid) {
		return
	}

	if self.Val < conf.P2 {
		body.(IPlayer).SendMsg(&msg.GS_TaskDailyValueChanged{
			Id:  confid,
			Val: self.Val,
		})
	} else {
		self.Fin = true

		body.(IPlayer).SendMsg(&msg.GS_TaskDailyItemCompleted{
			Id: confid,
		})

		// evt
		evtmgr.Fire(gconst.Evt_TaskDaily_Fin, body.(IPlayer), confid)
	}
}
