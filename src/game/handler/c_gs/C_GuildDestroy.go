package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildDestroy(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildDestroy)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildDestroy_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		if gld.LockMemberLeave() {
			return Err.Guild_LeaveTimeLimited
		}

		// get member
		m := gld.FindMember(plr.GetId())
		if m == nil {
			return Err.Guild_NotAMember
		}

		// check rank
		if m.Rank != guild.RK_Owner {
			return Err.Guild_LowPriv
		}

		// destroy
		return guild.GuildMgr.DestroyGuild(gld)
	}()

	plr.SendMsg(res)
}
