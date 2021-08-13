package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TaskMonthTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TaskMonthTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_TaskMonthTask_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetTaskMonth().Take(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.Id = req.Id

		return Err.OK
	}()

	plr.SendMsg(res)
}
