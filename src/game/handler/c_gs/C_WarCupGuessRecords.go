package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
)

func C_WarCupGuessRecords(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WarCupGuessRecords)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupGuessRecords_R{}

	res.Records = warcup.WarCupGuessRecords(plr)

	plr.SendMsg(res)
}
