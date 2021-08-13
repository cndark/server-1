package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_MiscGoldenHand(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_MiscGoldenHand)
	plr := ctx.(*app.Player)

	res := &msg.GS_MiscGoldenHand_R{}

	res.ErrorCode = func() int32 {
		ec, rwd := plr.GetMisc().GoldenHandClick()
		res.Rewards = rwd
		res.NextCrit = plr.GetMisc().GoldenHandCrit

		return ec
	}()

	plr.SendMsg(res)
}
