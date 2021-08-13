package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildInfoFull(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildInfoFull)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildInfoFull_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		// get info
		res.Info = gld.ToMsg_InfoFull(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
