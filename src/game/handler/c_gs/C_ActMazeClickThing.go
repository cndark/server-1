package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeClickThing(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMazeClickThing)

	plr := ctx.(*app.Player)

	res := &msg.GS_ActMazeClickThing_R{}
	res.ErrorCode = func() int32 {
		a := act.FindAct(gconst.ActName_Maze)
		if a == nil {
			return Err.Act_ActNotFound
		}

		if req.Pos < 0 || req.Pos >= 25 {
			return Err.Activity_MazePosLimit
		}

		ec, score, rwds := maze.ClickThing(plr, req.Pos, req.ItemId)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.Score = score

		return Err.OK
	}()

	plr.SendMsg(res)
}
