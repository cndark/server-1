package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WBossGetSelfRank(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WBossGetSelfRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_WBossGetSelfRank_R{}

	res.SelfRank, res.Jf = plr.GetWBoss().GetSelfRank()

	plr.SendMsg(res)
}
