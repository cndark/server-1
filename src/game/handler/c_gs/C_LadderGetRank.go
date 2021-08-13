package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_LadderGetRank(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_LadderGetRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_LadderGetRank_R{}

	res.Records = plr.GetLadder().GetLadderRank()

	plr.SendMsg(res)
}
