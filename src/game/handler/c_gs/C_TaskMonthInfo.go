package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_TaskMonthInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_TaskMonthInfo)

	plr := ctx.(*app.Player)

	res := &msg.GS_TaskMonthInfo_R{}

	res.Data = plr.GetTaskMonth().ToMsg()

	plr.SendMsg(res)
}
