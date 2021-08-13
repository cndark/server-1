package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_InviteTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_InviteTake)

	plr := ctx.(*app.Player)

	res := &msg.GS_InviteTake_R{}
	res.ErrorCode = func() int32 {

		if req.Id < 1 || req.Id > 3 {
			return Err.Failed
		}

		ec, rwds := plr.GetInvite().Take(req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
