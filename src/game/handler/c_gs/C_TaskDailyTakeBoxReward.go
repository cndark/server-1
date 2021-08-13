package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TaskDailyTakeBoxReward(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TaskDailyTakeBoxReward)
	plr := ctx.(*app.Player)

	res := &msg.GS_TaskDailyTakeBoxReward_R{}
	res.ErrorCode = func() int32 {
		er, rwds := plr.GetTaskDaily().TakeBoxReward(req.Id)
		if er != Err.OK {
			return er
		}

		res.Id = req.Id
		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
