package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildGetLog(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildGetLog)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildGetLog_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		res.Log = gld.ToMsg_Log()

		return Err.OK
	}()

	plr.SendMsg(res)
}
