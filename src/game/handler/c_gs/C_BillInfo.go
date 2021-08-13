package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_BillInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_BillInfo)
	plr := ctx.(*app.Player)

	res := plr.GetBill().ToMsg_Info()
	res.ErrorCode = Err.OK

	plr.SendMsg(res)
}
