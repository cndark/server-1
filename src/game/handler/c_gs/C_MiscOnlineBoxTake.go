package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MiscOnlineBoxTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MiscOnlineBoxTake)

	plr := ctx.(*app.Player)

	res := &msg.GS_MiscOnlineBoxTake_R{}

	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetMisc().TakeOnlineBox(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
