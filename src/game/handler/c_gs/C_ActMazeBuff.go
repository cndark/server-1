package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/maze"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMazeBuff(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActMazeBuff)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActMazeBuff_R{}

	a := act.FindAct(gconst.ActName_Maze)
	if a == nil {
		res.ErrorCode = Err.Act_ActNotFound
	}

	res.ErrorCode, res.BuffIds = maze.GetBuffIds(plr)

	plr.SendMsg(res)
}
