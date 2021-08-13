package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_GWarGetG2Members(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GWarGetG2Members)
	plr := ctx.(*app.Player)

	res := &msg.GS_GWarGetG2Members_R{}

	res.ErrorCode, res.Mbs = plr.GetGWar().GetG2Members()

	plr.SendMsg(res)
}
