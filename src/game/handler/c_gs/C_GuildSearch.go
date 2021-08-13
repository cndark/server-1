package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildSearch(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildSearch)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildSearch_R{}

	res.ErrorCode = func() int32 {
		gld := guild.GuildMgr.FindGuildByName(req.Name)
		if gld == nil {
			return Err.Guild_NotFound
		}

		// get info
		res.Row = gld.ToMsg_Row()

		return Err.OK
	}()

	plr.SendMsg(res)
}
