package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildWishItem(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildWishItem)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildWishItem_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		ec, seq := plr.GetGuildPlrData().WishItem(req.Num)
		res.Seq = seq
		return ec
	}()

	plr.SendMsg(res)
}
