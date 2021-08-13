package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MonthTicketTaskTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MonthTicketTaskTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_MonthTicketTaskTake_R{}
	res.ErrorCode = func() int32 {

		for i, id := range req.Ids {
			ec := plr.GetMonthTicket().TaskTake(id)
			if ec != Err.OK {
				if i == 0 {
					return ec
				} else {
					continue
				}
			}

			res.Ids = append(res.Ids, id)
		}

		res.Lv = plr.GetMonthTicket().Lv
		res.Exp = plr.GetMonthTicket().Exp

		return Err.OK
	}()

	plr.SendMsg(res)
}
