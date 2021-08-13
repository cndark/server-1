package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_TowerFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TowerFight)
	plr := ctx.(*app.Player)

	plr.GetTower().Fight(req.T, func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards) {
		res := &msg.GS_TowerFight_R{}

		res.ErrorCode = ec
		res.LvNum = plr.GetTower().LvNum
		res.Replay = replay
		res.Rewards = rwds

		plr.SendMsg(res)
	})
}
