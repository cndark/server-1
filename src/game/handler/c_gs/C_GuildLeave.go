package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildLeave(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildLeave)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildLeave_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		// get member
		m := gld.FindMember(plr.GetId())
		if m == nil {
			return Err.Guild_NotAMember
		}

		// leave
		if ec := m.Leave(guild.LR_Leave); ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
