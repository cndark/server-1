package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/arena"
	"fw/src/game/msg"
)

func C_ArenaRank(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArenaRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArenaRank_R{}
	res.SelfRank, res.Rows = arena.ArenaMgr.ToMsg_Rank(plr, req.Top)

	plr.SendMsg(res)

}
