package giftshop

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"strings"
)

// ============================================================================
// 礼包商店
type GiftShop struct {
	BuyCnt map[int32]int32

	plr IPlayer
}

// ============================================================================
func init() {

	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		prod_id := args[1].(int32)
		csext := args[2].(string)

		prod := gamedata.ConfBillProduct.Query(prod_id)
		if prod == nil || prod.TypeId != gconst.Bill_Gift {
			return
		}

		arr := strings.Split(csext, "_")
		if len(arr) < 2 || arr[0] != gconst.Bill_CsExt_Type_GiftShop {
			return
		}

		id := core.Atoi32(arr[1])
		plr.GetGiftShop().bill_buy(id, prod.PayId)

	})

	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		for id := range plr.GetGiftShop().BuyCnt {
			conf := gamedata.ConfGiftShop.Query(id)
			if conf == nil || conf.Reset == 1 {
				delete(plr.GetGiftShop().BuyCnt, id)
			}
		}
	})

	evtmgr.On(gconst.Evt_PlrResetWeekly, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		for id := range plr.GetGiftShop().BuyCnt {
			conf := gamedata.ConfGiftShop.Query(id)
			if conf == nil || conf.Reset == 2 {
				delete(plr.GetGiftShop().BuyCnt, id)
			}
		}
	})

	evtmgr.On(gconst.Evt_PlrResetMonthly, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		for id := range plr.GetGiftShop().BuyCnt {
			conf := gamedata.ConfGiftShop.Query(id)
			if conf == nil || conf.Reset == 3 {
				delete(plr.GetGiftShop().BuyCnt, id)
			}
		}
	})
}

func NewGiftShop() *GiftShop {
	return &GiftShop{
		BuyCnt: make(map[int32]int32),
	}
}

// ============================================================================

func (self *GiftShop) Init(plr IPlayer) {
	self.plr = plr
}

func (self *GiftShop) bill_buy(id int32, payid int32) {
	conf := gamedata.ConfGiftShop.Query(id)
	if conf == nil || conf.PayId != payid {
		return
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GiftShop)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.BuyCnt[id]++

	// res
	self.plr.SendMsg(&msg.GS_GiftShopNew{
		Id:      id,
		Rewards: rwds,
	})
}

func (self *GiftShop) Take(id int32) (int32, *msg.Rewards) {
	conf := gamedata.ConfGiftShop.Query(id)
	if conf == nil || conf.PayId != 0 {
		return Err.Failed, nil
	}

	if self.BuyCnt[id] >= conf.BuyCntLimit {
		return Err.Plr_TakenBefore, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GiftShop)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.BuyCnt[id]++

	return Err.OK, rwds
}

// ============================================================================

func (self *GiftShop) ToMsg() *msg.GiftShopData {
	return &msg.GiftShopData{
		BuyCnt: self.BuyCnt,
	}
}
