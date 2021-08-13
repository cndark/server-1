package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TargetDaysBuy(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TargetDaysBuy)
	plr := ctx.(*app.Player)

	res := &msg.GS_TargetDaysBuy_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetTargetDays().Buy(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
