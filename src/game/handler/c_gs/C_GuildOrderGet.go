package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildOrderGet(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildOrderGet)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildOrderGet_R{}

	res.ErrorCode = func() int32 {
		ec, records := plr.GetGuildPlrData().Order.GetOrders()
		res.Records = records
		res.GetOrdersTs = plr.GetGuildPlrData().Order.GetOrderTs.Unix()

		return ec
	}()

	plr.SendMsg(res)
}
