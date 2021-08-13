package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WBossGetRank(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WBossGetRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_WBossGetRank_R{}

	res.Rows = plr.GetWBoss().GetRank()

	plr.SendMsg(res)
}
