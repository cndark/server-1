package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/wboss"
	"fw/src/game/msg"
)

func C_WBossGetMaxDmgInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WBossGetMaxDmgInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_WBossGetMaxDmgInfo_R{}

	res.MaxDmgInfo = wboss.ToMsg_MaxDmgInfo()

	plr.SendMsg(res)
}
