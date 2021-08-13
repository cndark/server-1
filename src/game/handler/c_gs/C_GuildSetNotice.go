package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/core/wordsfilter"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildSetNotice(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildSetNotice)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildSetNotice_R{}

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

		// check notice
		if len(req.Notice) > C_guild_notice_maxlen {
			return Err.Guild_NoticeTooLong
		}

		if wordsfilter.IsSensitive(req.Notice) {
			return Err.Guild_InvalidNotice
		}

		// set
		if ec := gld.SetNotice(req.Notice); ec != Err.OK {
			return ec
		}

		// fire
		evtmgr.Fire(gconst.Evt_GuildChange, gld, plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
