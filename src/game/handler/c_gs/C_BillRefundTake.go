package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/refund"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_BillRefundTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_BillRefundTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_BillRefundTake_R{}

	res.ErrorCode = func() int32 {
		err, rwds := refund.TakeRefund(plr, req.Code)
		if err != Err.OK {
			return err
		}

		res.Rewards = rwds
		return Err.OK
	}()

	plr.SendMsg(res)
}
