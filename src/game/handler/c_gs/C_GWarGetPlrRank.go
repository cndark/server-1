package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GWarGetPlrRank(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GWarGetPlrRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_GWarGetPlrRank_R{}

	res.ErrorCode, res.Records = plr.GetGWar().GetPlrRank()

	plr.SendMsg(res)
}
