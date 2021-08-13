package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_LadderGetSummary(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_LadderGetSummary)
	plr := ctx.(*app.Player)

	res := &msg.GS_LadderGetSummary_R{}

	res.SelfRank = plr.GetLadder().GetSelfRank()

	plr.SendMsg(res)
}
