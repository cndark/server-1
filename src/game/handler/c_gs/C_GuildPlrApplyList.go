package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildPlrApplyList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildPlrApplyList)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildPlrApplyList_R{}

	res.ErrorCode = func() int32 {
		res.Rows = guild.GuildMgr.FetchGuildList_PlayerApplied(plr.GetId())
		return Err.OK
	}()

	plr.SendMsg(res)
}
