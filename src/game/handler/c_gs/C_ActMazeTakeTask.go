package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeTakeTask(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMazeTakeTask)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActMazeTakeTask_R{}

	res.ErrorCode = func() int32 {
		a := act.FindAct(gconst.ActName_Maze)
		if a == nil {
			return Err.Act_ActNotFound
		}

		ec, rwds := maze.TakeTask(plr, req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.Id = req.Id

		return Err.OK
	}()

	plr.SendMsg(res)
}
