package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_AppointRefresh(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_AppointRefresh)
	plr := ctx.(*app.Player)

	res := &msg.GS_AppointRefresh_R{}
	res.ErrorCode = func() int32 {

		ec := plr.GetAppoint().Refresh()
		if ec != Err.OK {
			return ec
		}

		res.Data = plr.GetAppoint().ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}
