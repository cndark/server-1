package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_FriendRemoveFrds(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_FriendRemoveFrds)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendRemoveFrds_R{}
	res.ErrorCode = func() int32 {

		ec := plr.GetFriend().Remove(req.PlrId)
		if ec != Err.OK {
			return ec
		}

		res.PlrId = req.PlrId

		return Err.OK
	}()

	plr.SendMsg(res)
}
