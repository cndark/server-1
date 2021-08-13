package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildOrderClose(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildOrderClose)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildOrderClose_R{}

	res.ErrorCode = func() int32 {
		ec, rwd := plr.GetGuildPlrData().Order.CloseOrder(req.Seq)
		res.Rewards = rwd
		return ec
	}()

	plr.SendMsg(res)
}
