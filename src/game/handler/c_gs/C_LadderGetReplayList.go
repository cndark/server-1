package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_LadderGetReplayList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_LadderGetReplayList)
	plr := ctx.(*app.Player)

	res := &msg.GS_LadderGetReplayList_R{}

	res.Records = plr.GetLadder().GetReplayList()

	plr.SendMsg(res)
}
