package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_LadderFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_LadderFight)
	plr := ctx.(*app.Player)

	res := &msg.GS_LadderFight_R{}

	plr.GetLadder().Fight(req.Team, req.TarId, req.TarRank, func(ec int32, replay *msg.BattleReplay, rwd *msg.Rewards) {
		res.ErrorCode = ec
		res.Replay = replay
		res.Rewards = rwd

		plr.SendMsg(res)
	})
}
