package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WLevelFundTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WLevelFundTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_WLevelFundTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetWLevelFund().Take(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
