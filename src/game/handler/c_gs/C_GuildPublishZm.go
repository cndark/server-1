package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/chat"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

func C_GuildPublishZm(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildPublishZm)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildPublishZm_R{}

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
		if m.Rank > guild.RK_Owner {
			return Err.Guild_LowPriv
		}

		// check cd
		now := time.Now()
		if now.Sub(gld.ZmTs).Hours() < 12 {
			return Err.Guild_ZmCd
		}

		// publish
		chat.Add(plr, &chat.ChatOne{
			Tp:     gconst.C_ChatType_GldZm,
			FromId: plr.GetId(),
			Ts:     now,
			GId:    gld.Id,
			GLv:    gld.Lv,
		})

		gld.ZmTs = now
		res.ZmTs = gld.ZmTs.Unix()

		return Err.OK
	}()

	plr.SendMsg(res)
}
