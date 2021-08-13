package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildBossFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildBossFight)
	plr := ctx.(*app.Player)

	plr.GetGuildPlrData().BossFight(req.Team, func(res *msg.GS_GuildBossFight_R) {
		plr.SendMsg(res)
	})
}
