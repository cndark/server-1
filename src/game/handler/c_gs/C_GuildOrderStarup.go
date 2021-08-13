package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildOrderStarup(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildOrderStarup)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildOrderStarup_R{}

	res.ErrorCode = func() int32 {
		return plr.GetGuildPlrData().Order.Starup(req.Seq)
	}()

	plr.SendMsg(res)
}
