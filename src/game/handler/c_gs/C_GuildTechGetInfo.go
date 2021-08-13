package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildTechGetInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildTechGetInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildTechGetInfo_R{}

	func() {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return
		}

		res.Techs = plr.GetGuildPlrData().Tech.Techs
	}()

	plr.SendMsg(res)
}
