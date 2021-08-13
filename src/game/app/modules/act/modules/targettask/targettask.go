package targettask

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

var actObj = &act_t{}

// ============================================================================

type act_t struct {
	act.BaseAct
}

type data_svr_t struct {
}

type data_plr_t struct {
	Attain  map[int32]*attain_obj_t // 统计进度
	Taken   []int32                 // 已经领取的
	ActGift map[int32]int32         // 活动礼包购买次数
}

type attain_obj_t struct {
	Id  int32
	Val float64 // progress value
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_TargetTask, actObj)

	// 更新购买活动礼包次数
	evtmgr.On(gconst.Evt_ActGift, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		actName := args[1].(string)
		id := args[2].(int32)

		if actObj.GetName() != actName {
			return
		}

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return
		}

		plr_data.ActGift[id]++
	})
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{
		Attain:  make(map[int32]*attain_obj_t),
		ActGift: make(map[int32]int32),
	}
}

func (self *act_t) ResetSvrData() {
}

func (self *act_t) ResetPlrData(iplr interface{}) {
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

func (self *act_t) OnQuit() {
}

func (self *act_t) OnStage() {
}

func (self *act_t) ToMsg(plr IPlayer) *msg.ActTargetTaskData {
	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	ret := &msg.ActTargetTaskData{
		Taken:   plr_data.Taken,
		ActGift: plr_data.ActGift,
	}

	for _, v := range plr_data.Attain {
		ret.Attain = append(ret.Attain, &msg.ActAttainObj{
			OId: v.Id,
			Val: v.Val,
		})
	}

	return ret
}

// ============================================================================

func (self *data_plr_t) get_attain_obj(oid int32) *attain_obj_t {
	obj := self.Attain[oid]
	if obj == nil {
		self.Attain[oid] = &attain_obj_t{Id: oid}
		return self.Attain[oid]
	}

	return obj
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
	if !actObj.Started() || !isChange {
		return
	}

	body.(IPlayer).SendMsg(&msg.GS_ActTargetTaskObjValueChanged{
		OId: self.Id,
		Val: self.Val,
	})
}

// ============================================================================

func ActTargetTaskInfo(plr IPlayer) *msg.ActTargetTaskData {
	return actObj.ToMsg(plr)
}

func TakeRewards(plr IPlayer, id int32) (int32, *msg.Rewards) {
	if !actObj.Started() && !actObj.Ended() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	for _, v := range plr_data.Taken {
		if v == id {
			return Err.Activity_TakeBefore, nil
		}
	}

	conf := gamedata.ConfActTargetTask.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	for _, v := range conf.Attain {
		obj := plr_data.get_attain_obj(v.AttainId)
		if obj.Val == 0 || obj.Val < v.P2 {
			return Err.Activity_CondLimited, nil
		}
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActTargetTask)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	plr_data.Taken = append(plr_data.Taken, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}
