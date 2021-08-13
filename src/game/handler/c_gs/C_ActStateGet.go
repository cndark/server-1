package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_ActStateGet(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActStateGet)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActStateGet_R{}

	res.Acts = plr.GetAct().ToMsg()

	plr.SendMsg(res)
}
