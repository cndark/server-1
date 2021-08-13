package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildWishClose(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildWishClose)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildWishClose_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		ec, rwd := plr.GetGuildPlrData().WishClose(req.Seq)
		res.Rewards = rwd
		return ec
	}()

	plr.SendMsg(res)
}
