package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GuildHarborDonateList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildHarborDonateList)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildHarborDonateList_R{}

	func() {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return
		}

		// list
		res.Records = gld.Harbor.ToMsg_DonateList()
	}()

	plr.SendMsg(res)
}
