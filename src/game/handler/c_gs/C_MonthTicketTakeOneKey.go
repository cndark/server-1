package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MonthTicketTakeOneKey(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_MonthTicketTakeOneKey)
	plr := ctx.(*app.Player)

	res := &msg.GS_MonthTicketTakeOneKey_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetMonthTicket().TakeOneKey()
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
