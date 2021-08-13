package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/billltday"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActBillLtDayInfo(message msg.Message, ctx interface{}) {
	plr := ctx.(*app.Player)

	res := &msg.GS_ActBillLtDayInfo_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_BillLtDay)
		if a == nil {
			return Err.Act_ActNotFound
		}

		res.Data = billltday.ActBillLtDayGetInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
