package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ShopGetInfo(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ShopGetInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_ShopGetInfo_R{}
	res.ErrorCode = func() int32 {
		if ec := plr.GetShop().IsShopOpen(req.ShopId); ec != Err.OK {
			return ec
		}

		plr.GetShop().Update(req.ShopId)

		ret := plr.GetShop().ToMsg(req.ShopId)
		if ret == nil {
			return Err.Failed
		}

		res.Shop = ret

		return Err.OK
	}()

	plr.SendMsg(res)
}
