package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActMazeInfo)

	res := &msg.GS_ActMazeInfo_R{}

	plr := ctx.(*app.Player)

	a := act.FindAct(gconst.ActName_Maze)
	if a == nil {
		res.ErrorCode = Err.Act_ActNotFound
	}

	res.Data = maze.ActMazeInfo(plr)

	plr.SendMsg(res)

}
