package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

func C_WarCupWatch(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WarCupWatch)
	plr := ctx.(*app.Player)

	evtmgr.Fire(gconst.Evt_WarCupWatch, plr)
}
