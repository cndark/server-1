package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/act/modules/monopoly"
	"fw/src/game/msg"
)

func C_ActMonopolyNPos(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMonopolyNPos)
	plr := ctx.(*app.Player)

	monopoly.NextNPos(plr, req.Step)

}
