package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GWarGetGuildRank(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GWarGetGuildRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_GWarGetGuildRank_R{}

	res.Records = plr.GetGWar().GetGuildRank()

	plr.SendMsg(res)
}
