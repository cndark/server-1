package wlevelfund

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 关卡基金
type WLevelFund struct {
	IsBuy bool    // 是否购买
	Taken []int32 // 领取id

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
			(conf_prod.PayId != gconst.Bill_PayId_WLevelFund) {
			return
		}

		plr.GetWLevelFund().IsBuy = true

		plr.SendMsg(&msg.GS_WLevelFundNew{
			FundId: conf_prod.PayId,
		})

		// fire
		evtmgr.Fire(gconst.Evt_BillWLevelFund, plr, conf_prod.PayId)
	})
}

func NewWLevelFund() *WLevelFund {
	return &WLevelFund{}
}

// ============================================================================

func (self *WLevelFund) Init(plr IPlayer) {
	self.plr = plr
}

func (self *WLevelFund) Take(id int32) (int32, *msg.Rewards) {
	if !self.IsBuy {
		return Err.Plr_NotPay, nil
	}

	for _, v := range self.Taken {
		if v == id {
			return Err.Plr_TakenBefore, nil
		}
	}

	conf := gamedata.ConfWorldLevelFund.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if self.plr.GetWLevelLvNum() < conf.Lv {
		return Err.WLevel_LowLvNum, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_WLevelFund)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	self.Taken = append(self.Taken, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// ============================================================================

func (self *WLevelFund) ToMsg() *msg.WLevelFundData {
	return &msg.WLevelFundData{
		IsBuy: self.IsBuy,
		Taken: self.Taken,
	}
}
