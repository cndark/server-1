package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_RiftMonsterFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RiftMonsterFight)
	plr := ctx.(*app.Player)

	plr.GetRift().FightMonster(req.T, func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards) {
		res := &msg.GS_RiftMonsterFight_R{}
		res.ErrorCode = ec
		res.Replay = replay
		res.Rewards = rwds

		plr.SendMsg(res)
	})
}
