package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_FriendApply(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_FriendApply)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendApply_R{}
	res.ErrorCode = func() int32 {

		ec := plr.GetFriend().Apply(req.PlrId)
		if ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
