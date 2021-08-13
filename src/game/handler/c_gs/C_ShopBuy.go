package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math"
)

func C_ShopBuy(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ShopBuy)
	plr := ctx.(*app.Player)

	res := &msg.GS_ShopBuy_R{}
	res.ErrorCode = func() int32 {
		if ec := plr.GetShop().IsShopOpen(req.ShopId); ec != Err.OK {
			return ec
		}

		// check bag limit
		if er := plr.GetBag().CheckFull(); er != Err.OK {
			return er
		}

		shopid, itemid, n := req.ShopId, req.ItemId, req.N
		if n < 1 || n > 999999 {
			n = 1
		}

		item := make(map[int32]int32)
		item[itemid] = n

		ec, buy_item, _, _ := shopid_buy(plr, shopid, item)
		if ec != Err.OK {
			return ec
		}

		res.ItemId = req.ItemId
		res.N = buy_item[req.ItemId]

		return Err.OK
	}()

	plr.SendMsg(res)
}

// 额外折扣
func shopid_extdiscount(plr *app.Player, shopid int32) float64 {
	ext_discount := float64(1)

	if shopid == 100 && plr.IsPrivCardValid(gconst.C_PrivCard_LongLife) {
		return 0.8
	}

	return ext_discount
}

// 购买一个商店
func shopid_buy(plr *app.Player, shopid int32, item map[int32]int32) (int32, map[int32]int32, map[int32]float64, map[int32]float64) {
	op := plr.GetBag().NewOp(gconst.ObjFrom_ShopBuy + shopid)

	buy_item := make(map[int32]int32)
	cost_item := make(map[int32]float64)
	rwds_item := make(map[int32]float64)

	for itemid, n := range item {
		//check
		er := plr.GetShop().CheckBuy(shopid, itemid)
		if er != Err.OK {
			return er, nil, nil, nil
		}

		conf := gamedata.ConfShopItem.Query(itemid)
		if conf == nil {
			return Err.Failed, nil, nil, nil
		}

		ec, price := utils.ConfBasePriceQuery(conf.Item, conf.Currency)
		if ec != Err.OK {
			return ec, nil, nil, nil
		}

		// modify n
		n = plr.GetShop().ModifyGoodBuyCnt(shopid, itemid, n)
		if n <= 0 {
			return Err.Shop_BuyCntTop, nil, nil, nil
		}

		//get price
		bn := float64(plr.GetShop().GetGoodBuyCnt(shopid, itemid))
		dis := float64(conf.Discount) / 10000 * shopid_extdiscount(plr, shopid)

		price = math.Floor(price * float64(conf.Num) * dis)
		price = float64(n)*(price+bn*conf.AddPrice) + conf.AddPrice*float64(n*(n-1))/2
		price = math.Floor(price)
		if price <= 0 {
			return Err.Shop_PriceError, nil, nil, nil
		}

		op.Dec(conf.Currency, price)
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec, nil, nil, nil
		}

		rwd := float64(conf.Num * n)
		op.Inc(conf.Item, rwd)

		buy_item[itemid] = n

		cost_item[conf.Currency] = price
		rwds_item[conf.Item] = rwd
	}

	for itemid, n := range buy_item {
		plr.GetShop().BuyGoods(shopid, itemid, n)
	}

	op.Apply()

	// fire
	evtmgr.Fire(gconst.Evt_ShopBuy, plr, shopid, buy_item)

	return Err.OK, buy_item, cost_item, rwds_item
}
