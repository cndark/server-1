package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MiscGldActGiftTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MiscGldActGiftTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_MiscGldActGiftTake_R{}
	res.ErrorCode = func() int32 {
		ec, rwds := plr.GetMisc().GldActGiftTake(req.Idx)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
