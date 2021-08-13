package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_DrawScoreBoxTake(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_DrawScoreBoxTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_DrawScoreBoxTake_R{}
	res.ErrorCode = func() int32 {

		err, rwds := plr.GetDraw().ScoreBoxTake()
		if err != Err.OK {
			return err
		}

		res.Score = plr.GetDraw().Score
		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
