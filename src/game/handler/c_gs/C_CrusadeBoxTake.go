package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_CrusadeBoxTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_CrusadeBoxTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_CrusadeBoxTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetCrusade().TakeBox(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Id = req.Id
		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
