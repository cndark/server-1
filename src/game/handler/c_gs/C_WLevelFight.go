package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WLevelFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WLevelFight)
	plr := ctx.(*app.Player)

	plr.GetWLevel().Fight(req.T, func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards) {
		res := &msg.GS_WLevelFight_R{}
		res.ErrorCode = ec
		res.Replay = replay
		res.Rewards = rwds
		res.LvNum = plr.GetWLevel().LvNum

		plr.SendMsg(res)
	})

}
