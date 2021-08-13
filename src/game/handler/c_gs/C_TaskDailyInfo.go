package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_TaskDailyInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_TaskDailyInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_TaskDailyInfo_R{}

	res.Data = plr.GetTaskDaily().ToMsg()

	plr.SendMsg(res)
}
