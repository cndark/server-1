package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildOrderList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildOrderList)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildOrderList_R{}

	res.Records = plr.GetGuildPlrData().Order.OrderList()
	res.GetOrdersTs = plr.GetGuildPlrData().Order.GetOrderTs.Unix()

	plr.SendMsg(res)
}
