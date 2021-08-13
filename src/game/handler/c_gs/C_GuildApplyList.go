package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildApplyList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildApplyList)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildApplyList_R{}

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

		// check rank
		if m.Rank > guild.RK_Vice {
			return Err.Guild_LowPriv
		}

		// get list
		res.Rows = gld.Apply.ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}
