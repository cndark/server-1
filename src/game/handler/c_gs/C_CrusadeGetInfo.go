package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_CrusadeGetInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_CrusadeGetInfo)
	plr := ctx.(*app.Player)

	res := plr.GetCrusade().GetInfo()

	plr.SendMsg(res)
}
