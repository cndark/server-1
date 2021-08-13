package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/monopoly"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMonopolyBuy(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMonopolyBuy)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActMonopolyBuy_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_Monopoly)
		if a == nil {
			return Err.Act_ActNotFound
		}

		ec, rwds := monopoly.Buy(plr, req.ShopIdx, req.ItemIdx)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
