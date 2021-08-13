package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_BillRefundCodeGet(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_BillRefundCodeGet)
	plr := ctx.(*app.Player)

	res := &msg.GS_BillRefundCodeGet_R{}

	res.ErrorCode = func() int32 {
		res.Code = plr.GetBill().GetRefundCode()
		return Err.OK
	}()

	plr.SendMsg(res)
}
