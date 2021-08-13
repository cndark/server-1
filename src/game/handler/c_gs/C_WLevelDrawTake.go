package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WLevelDrawTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WLevelDrawTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_WLevelDrawTake_R{}
	res.ErrorCode = func() int32 {
		ec, rwds, items := plr.GetWLevelDraw().Take(req.Idx, req.IsAutoDec)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.AutoDecItems = items

		return Err.OK
	}()

	plr.SendMsg(res)
}
