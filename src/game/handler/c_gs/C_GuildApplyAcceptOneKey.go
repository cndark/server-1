package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildApplyAcceptOneKey(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildApplyAcceptOneKey)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildApplyAcceptOneKey_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotAMember
		}

		// get member
		m := gld.FindMember(plr.GetId())
		if m == nil {
			return Err.Guild_NotAMember
		}

		// check rank
		if m.Rank > guild.RK_Vice {
			return Err.Guild_LowPriv
		}

		gld.ApplyAcceptOnekey()

		return Err.OK
	}()

	plr.SendMsg(res)
}
