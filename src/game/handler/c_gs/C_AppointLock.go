package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_AppointLock(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_AppointLock)
	plr := ctx.(*app.Player)

	res := &msg.GS_AppointLock_R{}
	res.ErrorCode = func() int32 {

		plr.GetAppoint().Lock(req.Seq, req.IsLock)

		res.Seq = req.Seq
		res.IsLock = req.IsLock

		return Err.OK
	}()

	plr.SendMsg(res)
}
