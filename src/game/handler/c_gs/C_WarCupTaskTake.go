package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupTaskTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WarCupTaskTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupTaskTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetWarCup().Take(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()
	plr.SendMsg(res)
}
