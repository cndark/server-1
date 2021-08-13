package growfund

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 成长基金
type GrowFund struct {
	IsBuy    bool    // 是否购买
	Taken    []int32 // 领取id
	TakenSvr []int32 // 领取服务器人数奖励

	plr IPlayer
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		pid := args[1].(int32)

		conf_prod := gamedata.ConfBillProduct.Query(pid)
		if conf_prod == nil {
			return
		}

		if (conf_prod.TypeId != gconst.Bill_Fund) ||
			(conf_prod.PayId != gconst.Bill_PayId_GrowFund) {
			return
		}

		plr.GetGrowFund().IsBuy = true

		plr.SendMsg(&msg.GS_GrowFundNew{
			FundId: conf_prod.PayId,
		})

		// fire
		evtmgr.Fire(gconst.Evt_BillGrowFund, plr, conf_prod.PayId)
	})
}

func NewGrowFund() *GrowFund {
	return &GrowFund{}
}

// ============================================================================

func (self *GrowFund) Init(plr IPlayer) {
	self.plr = plr
}

func (self *GrowFund) TakeLv(id int32) (int32, *msg.Rewards) {
	if !self.IsBuy {
		return Err.Plr_NotPay, nil
	}

	for _, v := range self.Taken {
		if v == id {
			return Err.Plr_TakenBefore, nil
		}
	}

	conf := gamedata.ConfGrowFund.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if self.plr.GetLevel() < conf.Lv {
		return Err.Plr_LowLevel, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GrowFund)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	self.Taken = append(self.Taken, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

func (self *GrowFund) TakeSvr(id int32) (int32, *msg.Rewards) {
	if !self.IsBuy {
		return Err.Plr_NotPay, nil
	}

	for _, v := range self.TakenSvr {
		if v == id {
			return Err.Plr_TakenBefore, nil
		}
	}

	conf := gamedata.ConfGrowFundReward.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if self.plr.GetLevel() < conf.Lv {
		return Err.Plr_LowLevel, nil
	}

	if GrowFundSvr.SvrCnt < conf.BuyCnt {
		return Err.Common_NeedMore, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GrowFund)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	self.TakenSvr = append(self.TakenSvr, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// ============================================================================

func (self *GrowFund) ToMsg() *msg.GrowFundData {
	return &msg.GrowFundData{
		IsBuy:     self.IsBuy,
		Taken:     self.Taken,
		TakenSvr:  self.TakenSvr,
		SvrBuyCnt: GrowFundSvr.SvrCnt,
	}
}
