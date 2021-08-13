package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WBossTakeMaxDmgRwd(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WBossTakeMaxDmgRwd)
	plr := ctx.(*app.Player)

	res := &msg.GS_WBossTakeMaxDmgRwd_R{}

	res.ErrorCode, res.Rewards = plr.GetWBoss().TakeMaxDmgRwd(req.N)

	plr.SendMsg(res)
}
