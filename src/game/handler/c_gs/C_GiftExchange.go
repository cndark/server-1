package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GiftExchange(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GiftExchange)
	plr := ctx.(*app.Player)

	res := &msg.GS_GiftExchange_R{}

	res.ErrorCode, res.Rewards = plr.GetMisc().GiftExchange(req.Code)

	plr.SendMsg(res)
}
