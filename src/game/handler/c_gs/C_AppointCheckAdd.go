package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_AppointCheckAdd(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_AppointCheckAdd)
	plr := ctx.(*app.Player)

	res := &msg.GS_AppointCheckAdd_R{}
	res.ErrorCode = func() int32 {

		plr.GetAppoint().CheckAdd()

		res.AddTs = plr.GetAppoint().AddTs.Unix()

		return Err.OK
	}()

	plr.SendMsg(res)
}
