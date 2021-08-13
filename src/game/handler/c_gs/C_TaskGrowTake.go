package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TaskGrowTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TaskGrowTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_TaskGrowTake_R{}
	res.ErrorCode = func() int32 {

		er, rwds := plr.GetTaskGrow().Take(req.Id)
		if er != Err.OK {
			return er
		}

		res.Rewards = rwds
		res.Id = req.Id

		return Err.OK
	}()

	plr.SendMsg(res)
}
