package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_FriendGiveAndRecv(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_FriendGiveAndRecv)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendGiveAndRecv_R{}
	res.ErrorCode = func() int32 {

		res.Cnt = plr.GetFriend().GiveAndRecv(req.PlrIds)

		res.PlrIds = req.PlrIds

		return Err.OK
	}()

	plr.SendMsg(res)
}
