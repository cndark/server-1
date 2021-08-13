package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_LadderGetReplay(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_LadderGetReplay)
	plr := ctx.(*app.Player)

	res := &msg.GS_LadderGetReplay_R{}

	plr.GetLadder().GetReplay(req.ReplayId, func(rp *msg.BattleReplay) {
		if rp == nil {
			res.ErrorCode = Err.Failed
		} else {
			res.ErrorCode = Err.OK
			res.Replay = rp
		}

		plr.SendMsg(res)
	})
}
