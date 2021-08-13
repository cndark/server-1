package bill

import (
	"fmt"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
)

// ============================================================================

const (
	C_refund_code_len = 12
)

// ============================================================================

type Bill struct {
	RefundCode   string          `bson:"refund_code"` // 返利码
	BuyCnt       map[int32]int32 // [payid]cnt			档位购买次数
	IsRealFirst  bool            // 是否充值过真实货币
	TotalBaseCcy int64           // 累计总充值基准货币

	plr IPlayer
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		// refund code generation & mail notification
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf != nil && conf.RefundSwitch == 1 {
			plr.GetBill().gen_refund_code()
		}
	})
}

// ============================================================================

func NewBill() *Bill {
	return &Bill{
		BuyCnt: make(map[int32]int32),

		RefundCode: "",
	}
}

func (self *Bill) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Bill) GiveItems(prod_id int32, csext string, amt int32, ccy string) bool {
	// give items
	prod := gamedata.ConfBillProduct.Query(prod_id)
	if prod == nil {
		log.Warning("bill-product NOT found:", prod_id, "user, amt", self.plr.GetId(), amt)
		return false
	}

	var rwds *msg.Rewards
	switch prod.TypeId {
	case gconst.Bill_Normal:
		rwds = self.type_normal(prod_id)

	case gconst.Bill_PrivCard:
		if (prod.PayId == gconst.Bill_PayId_SaleWeek ||
			prod.PayId == gconst.Bill_PayId_SaleMonth) &&
			self.BuyCnt[prod.PayId] > 0 {
			return false
		}
		rwds = self.type_card(prod_id)

	case gconst.Bill_First:
		rwds = self.type_first(prod_id)

	case gconst.Bill_Gift:
		rwds = self.type_gift(prod_id)

	case gconst.Bill_Fund, gconst.Bill_RushCross, gconst.Bill_MonthTicket:
		rwds = self.type_default(prod_id)

	default:
		log.Warning("bill-product typeid error:", prod.TypeId, "user, amt", self.plr.GetId(), amt)
		return false
	}

	// send notification mail
	self.send_mail_for_prod(prod_id)

	// update total info
	self.update_total_info(prod.BaseCcy)

	//buy cnt
	self.BuyCnt[prod.PayId]++

	// notify
	self.plr.SendMsg(&msg.GS_BillDone{
		ProdId:  prod_id,
		Rewards: rwds,
	})

	// fire
	evtmgr.Fire(gconst.Evt_BillDone, self.plr, prod_id, csext)

	return true
}

func (self *Bill) type_normal(prod_id int32) (rwds *msg.Rewards) {
	prod := gamedata.ConfBillProduct.Query(prod_id)

	// deliver
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_Bill_Normal)

	// fixed items
	for _, v := range prod.Goods {
		op.Inc(v.Id, int64(v.N))
	}

	// ext items
	cnt := self.BuyCnt[prod.PayId] // by price
	if cnt == 0 {
		// first-pay
		for _, v := range prod.ExtGoods1 {
			op.Inc(v.Id, int64(v.N))
		}
	} else {
		// normal-pay
		for _, v := range prod.ExtGoods2 {
			op.Inc(v.Id, int64(v.N))
		}
	}

	// apply
	rwds = op.Apply().ToMsg()

	return
}

func (self *Bill) type_card(prod_id int32) (rwds *msg.Rewards) {
	prod := gamedata.ConfBillProduct.Query(prod_id)

	// deliver
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_Bill_Card)

	// fixed items
	for _, v := range prod.Goods {
		op.Inc(v.Id, int64(v.N))
	}

	// apply
	rwds = op.Apply().ToMsg()

	return
}

func (self *Bill) type_first(prod_id int32) (rwds *msg.Rewards) {
	prod := gamedata.ConfBillProduct.Query(prod_id)

	// deliver
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_Bill_First)

	// fixed items
	for _, v := range prod.Goods {
		op.Inc(v.Id, int64(v.N))
	}

	// apply
	rwds = op.Apply().ToMsg()

	return
}

func (self *Bill) type_default(prod_id int32) (rwds *msg.Rewards) {
	// do nothing
	return
}

func (self *Bill) type_gift(prod_id int32) (rwds *msg.Rewards) {
	prod := gamedata.ConfBillProduct.Query(prod_id)

	// deliver
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_Bill_Gift)

	// fixed items
	for _, v := range prod.Goods {
		op.Inc(v.Id, int64(v.N))
	}

	// apply
	rwds = op.Apply().ToMsg()

	return
}

func (self *Bill) send_mail_for_prod(prod_id int32) {
	prod := gamedata.ConfBillProduct.Query(prod_id)
	if prod.MailId == 0 {
		return
	}

	conf_mail := gamedata.ConfMail.Query(prod.MailId)
	if conf_mail == nil {
		return
	}

	mail.New(self.plr).SetKey(conf_mail.Key).Send()
}

func (self *Bill) update_total_info(v int32) {
	// update
	self.TotalBaseCcy += int64(v)

	// flush
	async.Push(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{"$set": db.M{"base.bill.totalbaseccy": self.TotalBaseCcy}},
		)
		if err != nil {
			log.Error("update_total_info() failed:", err)
		}
	})
}

func (self *Bill) GetTotalString() string {
	return fmt.Sprintf("%d", self.TotalBaseCcy)
}

func (self *Bill) gen_refund_code() {
	if self.RefundCode != "" {
		return
	}

	// gen code
	code := utils.GenRandCode(C_refund_code_len)
	self.RefundCode = code

	// flush
	async.Push(
		func() {
			err := dbmgr.DBBill.Insert(
				"refund",
				db.M{
					"_id":  self.plr.GetId(),
					"code": code,
				},
			)
			if err != nil {
				log.Error("flush user refund code failed:", err)
			}
		},
	)

	// send mail
	conf_mail := gamedata.ConfMailM.Query("refund")
	if conf_mail != nil {
		mail.New(self.plr).SetKey(conf_mail.Key).AddDict("refundCode", code).Send()
	}
}

func (self *Bill) GetRefundCode() string {
	return self.RefundCode
}

// ============================================================================

func (self *Bill) ToMsg_Info() *msg.GS_BillInfo_R {
	m := &msg.GS_BillInfo_R{}

	for id, cnt := range self.BuyCnt {
		m.BuyCnt = append(m.BuyCnt, &msg.BillBuyCnt{id, cnt})
	}

	m.TotalBaseCcy = self.TotalBaseCcy

	return m
}
