package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_CounterBuy(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_CounterBuy)
	plr := ctx.(*app.Player)

	res := &msg.GS_CounterBuy_R{}
	res.ErrorCode = func() int32 {

		ec := plr.GetCounter().Buy(req.Id)
		if ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
