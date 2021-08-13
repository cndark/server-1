package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_CrusadeFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_CrusadeFight)
	plr := ctx.(*app.Player)
	// fight
	plr.GetCrusade().Fight(req.T,
		func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards, hploss map[int64]float64) {
			res := &msg.GS_CrusadeFight_R{}

			res.ErrorCode = ec
			res.Replay = replay
			res.Rewards = rwds
			res.HpLoss = hploss

			plr.SendMsg(res)
		})
}
