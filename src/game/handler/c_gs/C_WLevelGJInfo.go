package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WLevelGJInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WLevelGJInfo)
	plr := ctx.(*app.Player)

	plr.GetWLevel().GJLoot()

	res := &msg.GS_WLevelGJInfo_R{
		GJLootRwd: plr.GetWLevel().GJLootRwd,
	}

	plr.SendMsg(res)
}
