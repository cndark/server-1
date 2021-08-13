package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeClickBattle(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMazeClickBattle)

	plr := ctx.(*app.Player)

	res := &msg.GS_ActMazeClickBattle_R{}

	a := act.FindAct(gconst.ActName_Maze)
	if a == nil {
		res.ErrorCode = Err.Act_ActNotFound
	}

	if req.Pos < 0 || req.Pos >= 25 {
		res.ErrorCode = Err.Activity_MazePosLimit
	}

	maze.ClickBattle(plr, req.Pos, req.T, func(ec, score int32, replay *msg.BattleReplay, rwds *msg.Rewards) {

		res.ErrorCode = ec
		res.Replay = replay
		res.Rewards = rwds
		res.Score = score

		plr.SendMsg(res)
	})

}
