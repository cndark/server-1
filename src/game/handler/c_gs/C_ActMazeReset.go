package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeReset(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActMazeReset)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActMazeReset_R{}

	res.ErrorCode = func() int32 {
		a := act.FindAct(gconst.ActName_Maze)
		if a == nil {
			return Err.Act_ActNotFound
		}

		ec, data := maze.Reset(plr)
		if ec != Err.OK {
			return ec
		}

		res.Data = data

		return Err.OK
	}()

	plr.SendMsg(res)
}
