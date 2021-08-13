package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_MonthTicketInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_MonthTicketInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_MonthTicketInfo_R{}
	res.Data = plr.GetMonthTicket().ToMsg()

	plr.SendMsg(res)
}
