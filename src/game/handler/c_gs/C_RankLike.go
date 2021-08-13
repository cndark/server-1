package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_RankLike(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RankLike)
	plr := ctx.(*app.Player)

	res := &msg.GS_RankLike_R{}

	res.ErrorCode = func() int32 {
		ec, rwd := plr.GetRankPlay().Like(req.RkId, req.PlrId)
		res.Rewards = rwd

		return ec
	}()

	plr.SendMsg(res)
}
