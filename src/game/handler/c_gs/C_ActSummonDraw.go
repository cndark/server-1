package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/act/modules/summon"
	"fw/src/game/msg"
)

func C_ActSummonDraw(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActSummonDraw)
	plr := ctx.(*app.Player)

	summon.Draw(plr, req.IsDiam, req.N)
}
