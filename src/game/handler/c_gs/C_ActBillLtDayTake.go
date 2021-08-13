package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/billltday"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActBillLtDayTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActBillLtDayTake)

	plr := ctx.(*app.Player)

	res := &msg.GS_ActBillLtDayTake_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_BillLtDay)
		if a == nil {
			return Err.Act_ActNotFound
		}

		ec, rwds := billltday.Take(plr, req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
