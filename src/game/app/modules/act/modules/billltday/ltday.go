package billltday

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
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
	Ts       time.Time // 上次增加时间
	Day      int32     // 充值天数
	DailyNum int32     // 当天充值数
	Taken    []int32
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_BillLtDay, actObj)

	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return
		}

		L := len(gamedata.ConfActBillLtDayM.QueryItems(actObj.GetConfGrp()))
		if len(plr_data.Taken) >= L {
			plr_data.Day = 0
			plr_data.Taken = []int32{}
		}

		plr_data.DailyNum = 0
	})

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

		conf_g := gamedata.ConfActivityPublic.Query(1)
		if conf_g == nil {
			return
		}

		now := time.Now()
		if !core.IsSameDay(now, plr_data.Ts) && (plr_data.DailyNum < conf_g.ActBillLtDayMin) &&
			(plr_data.DailyNum+conf_prod.BaseCcy >= conf_g.ActBillLtDayMin) {
			plr_data.Day++
			plr_data.Ts = now
		}

		plr_data.DailyNum += conf_prod.BaseCcy
	})
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{}
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

func (self *act_t) OnStage() {
}

// ============================================================================

func (self *act_t) ToMsg(plr IPlayer) *msg.ActBillLtDayData {
	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	ret := &msg.ActBillLtDayData{
		BillDay: plr_data.Day,
		Taken:   plr_data.Taken,
	}

	return ret
}

// ============================================================================

func ActBillLtDayGetInfo(plr IPlayer) *msg.ActBillLtDayData {
	return actObj.ToMsg(plr)
}

func Take(plr IPlayer, id int32) (int32, *msg.Rewards) {
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

	conf := gamedata.ConfActBillLtDay.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if conf.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	if plr_data.Day < conf.BillDay {
		return Err.Activity_BillDayNotEnough, nil
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActBillLtDay)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	plr_data.Taken = append(plr_data.Taken, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}
