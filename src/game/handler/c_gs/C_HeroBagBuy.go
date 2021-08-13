package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroBagBuy(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_HeroBagBuy)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroBagBuy_R{}
	res.ErrorCode = func() int32 {

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_HeroBagBuy)
		cop.DecCounter(gconst.Cnt_HeroBagBuy, 1)
		if ec := cop.CheckEnough(); ec != Err.OK {
			return ec
		}

		cop.Apply()

		return Err.OK
	}()

	plr.SendMsg(res)
}
