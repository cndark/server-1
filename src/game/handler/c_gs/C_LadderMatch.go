package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_LadderMatch(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_LadderMatch)
	plr := ctx.(*app.Player)

	res := &msg.GS_LadderMatch_R{}

	plr.GetLadder().Match(func(ec int32, ret []*msg.LadderPlayerInfo) {
		res.ErrorCode = ec
		res.Plrs = ret
		plr.SendMsg(res)
	})
}
