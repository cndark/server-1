package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/billlttotal"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActBillLtTotalInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActBillLtTotalInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActBillLtTotalInfo_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_BillLtTotal)
		if a == nil {
			return Err.Act_ActNotFound
		}

		res.Data = billlttotal.ActBillLtTotalGetInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
