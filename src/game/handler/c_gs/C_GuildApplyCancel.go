package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildApplyCancel(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildApplyCancel)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildApplyCancel_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := guild.GuildMgr.FindGuild(req.GuildId)
		if gld == nil {
			return Err.Guild_NotFound
		}

		// apply
		if ec := gld.ApplyCancel(plr); ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
