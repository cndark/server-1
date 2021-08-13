package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/core/wordsfilter"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"

	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"strings"
	"time"
)

const (
	C_guild_name_maxlen   = 24
	C_guild_notice_maxlen = 240
)

func C_GuildCreate(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildCreate)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildCreate_R{}

	res.ErrorCode = func() int32 {
		// has guild ?
		if plr.GetGuild() != nil {
			return Err.Guild_PlrInGuild
		}

		// check name
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			return Err.Guild_InvalidName
		} else if len(req.Name) > C_guild_name_maxlen {
			return Err.Guild_NameTooLong
		}

		if wordsfilter.IsSensitive(req.Name) {
			return Err.Guild_InvalidName
		}

		// check notice
		if len(req.Notice) > C_guild_notice_maxlen {
			return Err.Guild_NoticeTooLong
		}

		if wordsfilter.IsSensitive(req.Notice) {
			return Err.Guild_InvalidNotice
		}

		// conf
		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return Err.Failed
		}

		// check join cd
		if time.Now().Before(plr.GetMisc().GldLeaveTs) {
			return Err.Guild_JoinTimeLimited
		}

		// cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_GuildCreate)
		for _, v := range conf_g.GuildCreateCost {
			op.Dec(v.Id, int64(v.N))
		}
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// create
		ec, gld := guild.GuildMgr.CreateGuild(plr, req.Name, req.Notice, req.Icon)
		if ec != Err.OK {
			return ec
		}

		// apply
		op.Apply()

		// fire
		evtmgr.Fire(gconst.Evt_GuildChange, gld, plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
