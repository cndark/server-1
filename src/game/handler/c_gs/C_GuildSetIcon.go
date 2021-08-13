package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildSetIcon(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildSetIcon)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildSetIcon_R{}

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

		// set
		if ec := gld.SetIcon(req.Icon); ec != Err.OK {
			return ec
		}

		// fire
		evtmgr.Fire(gconst.Evt_GuildChange, gld, plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
