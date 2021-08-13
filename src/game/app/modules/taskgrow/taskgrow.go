package taskgrow

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================
// 进阶之路
type TaskGrow struct {
	Taken []int32

	plr IPlayer
}

// ============================================================================

func NewTaskGrow() *TaskGrow {
	return &TaskGrow{}
}

// ============================================================================

func (self *TaskGrow) Init(plr IPlayer) {
	self.plr = plr
}

func (self *TaskGrow) IsComplete(id int32) bool {
	for _, v := range self.Taken {
		if v == id {
			return true
		}
	}

	return false
}

func (self *TaskGrow) Take(id int32) (int32, *msg.Rewards) {
	if self.IsComplete(id) {
		return Err.TaskGrow_AlreadRewarded, nil
	}

	conf := gamedata.ConfTaskGrow.Query(id)
	if conf == nil || len(conf.AttainTab) == 0 {
		return Err.Failed, nil
	}

	for _, v := range conf.AttainTab {
		val := self.plr.GetAttainObjVal(v.AttainId)

		if val == 0 || val < v.P2 {
			return Err.TaskGrow_NotCompleted, nil
		}
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TaskGrow)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.Taken = append(self.Taken, id)

	// evt
	evtmgr.Fire(gconst.Evt_TaskGrow_Take, self.plr, id)

	return Err.OK, rwds
}

// ============================================================================

func (self *TaskGrow) ToMsg() *msg.TaskGrowData {
	return &msg.TaskGrowData{
		Taken: self.Taken,
	}
}
