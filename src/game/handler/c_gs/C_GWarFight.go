package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GWarFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GWarFight)
	plr := ctx.(*app.Player)

	plr.GetGWar().Fight(req.TarId, req.Team, func(res *msg.GS_GWarFight_R) {
		plr.SendMsg(res)
	})
}
