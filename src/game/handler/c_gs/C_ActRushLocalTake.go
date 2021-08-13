package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/act/modules/rushlocal"
	"fw/src/game/msg"
)

func C_ActRushLocalTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActRushLocalTake)
	plr := ctx.(*app.Player)

	rushlocal.TakeRankRewards(plr, req.ActName, req.RankId)
}
