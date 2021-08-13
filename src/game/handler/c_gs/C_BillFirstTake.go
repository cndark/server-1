package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_BillFirstTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_BillFirstTake)

	plr := ctx.(*app.Player)

	res := &msg.GS_BillFirstTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetBillFirst().Take(req.Id, req.Day)

		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)

}
