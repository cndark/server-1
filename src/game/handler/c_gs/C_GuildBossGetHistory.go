package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildBossGetHistory(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildBossGetHistory)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildBossGetHistory_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotAMember
		}

		// get
		ec, info := gld.Boss.ToMsg_History(req.Num)
		res.Info = info

		return ec
	}()

	plr.SendMsg(res)
}
