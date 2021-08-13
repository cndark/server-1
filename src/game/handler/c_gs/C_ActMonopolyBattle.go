package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/monopoly"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMonopolyBattle(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActMonopolyBattle)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActMonopolyBattle_R{}

	if act.FindAct(gconst.ActName_Monopoly) == nil {
		res.ErrorCode = Err.Act_ActNotFound
	}

	monopoly.Battle(plr, req.Tp, req.Idx, req.T, func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards) {

		res.ErrorCode = ec
		res.Replay = replay
		res.Rewards = rwds

		plr.SendMsg(res)
	})

}
