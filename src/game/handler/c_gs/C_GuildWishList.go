package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildWishList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildWishList)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildWishList_R{}

	res.Wishes = plr.GetGuildPlrData().WishList()

	plr.SendMsg(res)
}
