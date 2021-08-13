package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_CounterRecover(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_CounterRecover)
	plr := ctx.(*app.Player)

	res := &msg.GS_CounterRecover_R{}
	res.ErrorCode = func() int32 {

		ec, cnt, ts := plr.GetCounter().CheckRecover(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Id = req.Id
		res.Cnt = cnt
		res.Ts = ts

		return Err.OK
	}()

	plr.SendMsg(res)
}
