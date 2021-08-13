package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildList(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildList)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildList_R{}

	res.ErrorCode = func() int32 {
		res.Rows = guild.GuildMgr.FetchGuildList(req.Page)
		return Err.OK
	}()

	plr.SendMsg(res)
}
