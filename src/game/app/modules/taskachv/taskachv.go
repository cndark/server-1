package taskachv

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 成就任务
type TaskAchv struct {
	Taken []int32 // 完成成就任务

	plr IPlayer
}

// ============================================================================

func NewTaskAchv() *TaskAchv {
	return &TaskAchv{}
}

// ============================================================================

func (self *TaskAchv) Init(plr IPlayer) {
	self.plr = plr
}

func (self *TaskAchv) IsComplete(id int32) bool {
	for _, v := range self.Taken {
		if v == id {
			return true
		}
	}

	return false
}

func (self *TaskAchv) Take(id int32) (int32, *msg.Rewards) {
	if self.IsComplete(id) {
		return Err.TaskAchv_AlreadRewarded, nil
	}

	conf := gamedata.ConfTaskAchv.Query(id)
	if conf == nil || len(conf.AttainTab) == 0 {
		return Err.Failed, nil
	}

	for _, v := range conf.AttainTab {
		val := self.plr.GetAttainObjVal(v.AttainId)

		if val == 0 || val < v.P2 {
			return Err.TaskAchv_NotCompleted, nil
		}
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_TaskAchv)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.Taken = append(self.Taken, id)

	// evt
	evtmgr.Fire(gconst.Evt_TaskAchv_Take, self.plr, id)

	return Err.OK, rwds
}

// ============================================================================

func (self *TaskAchv) ToMsg() *msg.TaskAchvData {
	return &msg.TaskAchvData{
		Taken: self.Taken,
	}
}
