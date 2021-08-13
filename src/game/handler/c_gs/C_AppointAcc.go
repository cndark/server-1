package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_AppointAcc(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_AppointAcc)
	plr := ctx.(*app.Player)

	res := &msg.GS_AppointAcc_R{}
	res.ErrorCode = func() int32 {

		ec := plr.GetAppoint().Acc(req.Seq)
		if ec != Err.OK {
			return ec
		}

		res.Seq = req.Seq

		return Err.OK
	}()

	plr.SendMsg(res)
}
