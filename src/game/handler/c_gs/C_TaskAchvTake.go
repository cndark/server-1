package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TaskAchvTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TaskAchvTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_TaskAchvTake_R{}
	res.ErrorCode = func() int32 {

		er, rwds := plr.GetTaskAchv().Take(req.Id)
		if er != Err.OK {
			return er
		}

		res.Rewards = rwds
		res.Id = req.Id

		return Err.OK
	}()

	plr.SendMsg(res)
}
