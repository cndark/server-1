package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GWarGetSummary(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GWarGetSummary)
	plr := ctx.(*app.Player)

	res := plr.GetGWar().GetSummary()

	plr.SendMsg(res)
}
