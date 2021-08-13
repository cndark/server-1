package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WBossGetSummary(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WBossGetSummary)
	plr := ctx.(*app.Player)

	res := plr.GetWBoss().ToMsg_Summary()

	plr.SendMsg(res)
}
