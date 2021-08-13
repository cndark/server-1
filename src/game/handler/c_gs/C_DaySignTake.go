package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_DaySignTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_DaySignTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_DaySignTake_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetDaySign().Take(req.Id)
		if ec != Err.OK {
			return ec
		}
		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
