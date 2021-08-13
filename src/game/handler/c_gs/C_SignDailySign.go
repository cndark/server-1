package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_SignDailySign(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_SignDailySign)
	plr := ctx.(*app.Player)

	res := &msg.GS_SignDailySign_R{}
	res.ErrorCode = func() int32 {
		ec, rwds := plr.GetSignDaily().Sign()
		if ec != Err.OK {
			return ec
		}
		res.Rewards = rwds
		return Err.OK
	}()

	plr.SendMsg(res)
}
