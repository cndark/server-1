package billlttotal

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

type data_plr_t struct {
	Total int64
	Taken []int32
}

type data_svr_t struct {
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_BillLtTotal, actObj)

	on_bill_done()
}

func on_bill_done() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		if !actObj.Started() {
			return
		}

		plr := args[0].(IPlayer)
		pid := args[1].(int32)

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return
		}

		conf_prod := gamedata.ConfBillProduct.Query(pid)
		if conf_prod == nil {
			return
		}

		plr_data.Total += int64(conf_prod.BaseCcy)

		plr.SendMsg(&msg.GS_ActBillLtTotal{
			Total: plr_data.Total,
		})
	})
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{}
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

func (self *act_t) OnStage() {
}

// ============================================================================

func (self *act_t) ToMsg(plr IPlayer) *msg.ActBillLtTotalData {
	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	ret := &msg.ActBillLtTotalData{
		Total: plr_data.Total,
		Taken: plr_data.Taken,
	}

	return ret
}

// ============================================================================

func ActBillLtTotalGetInfo(plr IPlayer) *msg.ActBillLtTotalData {
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

	conf := gamedata.ConfActBillLtTotal.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if conf.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	if plr_data.Total < int64(conf.BillCond) {
		return Err.Activity_BillNotEnough, nil
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActBillLtTotal)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	plr_data.Taken = append(plr_data.Taken, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}
