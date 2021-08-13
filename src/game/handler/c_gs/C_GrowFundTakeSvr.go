package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GrowFundTakeSvr(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GrowFundTakeSvr)
	plr := ctx.(*app.Player)

	res := &msg.GS_GrowFundTakeSvr_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetGrowFund().TakeSvr(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
