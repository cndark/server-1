package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ShopRefresh(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ShopRefresh)
	plr := ctx.(*app.Player)

	res := &msg.GS_ShopRefresh_R{}
	res.ErrorCode = func() int32 {
		conf := gamedata.ConfShop.Query(req.ShopId)
		if conf == nil {
			return Err.Failed
		}

		if ec := plr.GetShop().IsShopOpen(req.ShopId); ec != Err.OK {
			return ec
		}

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_ShopRefresh)
		if conf.FreeCounter > 0 && plr.GetCounter().GetRemain(conf.FreeCounter) > 0 {
			cop.DecCounter(conf.FreeCounter, 1)
		} else {
			if len(conf.RefreshCost) == 0 {
				return Err.Shop_RefreshCntFull
			}

			for _, v := range conf.RefreshCost {
				cop.Dec(v.Id, v.N)
			}

			if ec := cop.CheckEnough(); ec != Err.OK {
				return ec
			}
		}

		if ec := plr.GetShop().Refresh(req.ShopId); ec != Err.OK {
			return ec
		}

		cop.Apply()

		res.Shop = plr.GetShop().ToMsg(req.ShopId)

		return Err.OK
	}()

	plr.SendMsg(res)
}
