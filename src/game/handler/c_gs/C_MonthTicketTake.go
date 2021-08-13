package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MonthTicketTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MonthTicketTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_MonthTicketTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetMonthTicket().Take(req.Lv)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.TakeBase = plr.GetMonthTicket().TakeBase
		res.TakeTicket = plr.GetMonthTicket().TakeTicket

		return Err.OK
	}()

	plr.SendMsg(res)
}
