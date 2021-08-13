package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MonthTicketBuyUp(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MonthTicketBuyUp)
	plr := ctx.(*app.Player)

	res := &msg.GS_MonthTicketBuyUp_R{}
	res.ErrorCode = func() int32 {
		if req.N <= 0 || req.N > 50 {
			return Err.Failed
		}

		ec := plr.GetMonthTicket().BuyUp(req.N)
		if ec != Err.OK {
			return ec
		}

		res.Lv = plr.GetMonthTicket().Lv
		res.Exp = plr.GetMonthTicket().Exp

		return Err.OK
	}()

	plr.SendMsg(res)
}
