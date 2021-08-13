package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/act/modules/msummon"
	"fw/src/game/msg"
)

func C_ActMSummonDraw(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMSummonDraw)
	plr := ctx.(*app.Player)

	msummon.Draw(plr, req.IsDiam, req.N)
}
