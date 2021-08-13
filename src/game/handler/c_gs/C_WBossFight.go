package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WBossFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WBossFight)
	plr := ctx.(*app.Player)

	plr.GetWBoss().Fight(req.T, func(r *msg.GS_WBossFight_R) {
		plr.SendMsg(r)
	})
}
