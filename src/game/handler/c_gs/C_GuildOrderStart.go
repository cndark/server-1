package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildOrderStart(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildOrderStart)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildOrderStart_R{}

	res.ErrorCode = func() int32 {
		ec, ts := plr.GetGuildPlrData().Order.StartOrder(req.Seq)
		res.StartTs = ts
		return ec
	}()

	plr.SendMsg(res)
}
