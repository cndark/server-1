package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeClick(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMazeClick)

	plr := ctx.(*app.Player)

	res := &msg.GS_ActMazeClick_R{}
	a := act.FindAct(gconst.ActName_Maze)
	if a == nil {
		res.ErrorCode = Err.Act_ActNotFound
	}

	if req.Pos < 0 || req.Pos >= 25 {
		res.ErrorCode = Err.Activity_MazePosLimit
	}

	maze.Click(plr, req.Pos, func(ec, score, r_seq int32, bd *msg.BattleData, rwds *msg.Rewards) {

		res.ErrorCode = ec
		res.Battle = bd
		res.Rewards = rwds
		res.Score = score
		res.Seq = r_seq

	})

	plr.SendMsg(res)

}
