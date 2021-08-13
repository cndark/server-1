package c_gs

import (
	"fw/src/core/wordsfilter"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"strings"
)

func C_GuildChangeName(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildChangeName)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildChangeName_R{}

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

		// check name
		req.Name = strings.TrimSpace(req.Name)

		if req.Name == "" {
			return Err.Guild_InvalidName
		} else if req.Name == gld.Name {
			return Err.Guild_SameName
		} else if len(req.Name) > C_guild_name_maxlen {
			return Err.Guild_NameTooLong
		}

		if wordsfilter.IsSensitive(req.Name) {
			return Err.Guild_InvalidName
		}

		// check cost
		conf := gamedata.ConfGlobalPublic.Query(1)
		op := plr.GetBag().NewOp(gconst.ObjFrom_GuildChangeName)

		for _, v := range conf.GuildNameChange {
			op.Dec(v.Id, int64(v.N))
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// change name
		gld.ChangeName(req.Name, func(b bool) {
			if b {
				op.Apply()

				res.ErrorCode = Err.OK
			} else {
				res.ErrorCode = Err.Guild_DupName
			}

			plr.SendMsg(res)
		})

		return Err.OK
	}()

	if res.ErrorCode != Err.OK {
		plr.SendMsg(res)
	}
}
