package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_AppointTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_AppointTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_AppointTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetAppoint().Take(req.Seq)
		if ec != Err.OK {
			return ec
		}

		res.Seq = req.Seq
		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
