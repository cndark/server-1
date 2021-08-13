package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_WLevelFightOneKey(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WLevelFightOneKey)
	plr := ctx.(*app.Player)

	plr.GetWLevel().FightOneKey(req.T, func(ec int32, replays []*msg.BattleReplay, rwds *msg.Rewards) {
		res := &msg.GS_WLevelFightOneKey_R{}
		res.ErrorCode = ec
		res.Replay = replays
		res.Rewards = rwds
		res.LvNum = plr.GetWLevel().LvNum

		plr.SendMsg(res)
	})
}
