package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_ArenaRecordInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ArenaRecordInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArenaRecordInfo_R{}
	res.Records = plr.GetArena().Records_ToMsg()

	plr.SendMsg(res)
}
