package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_FriendAccept(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_FriendAccept)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendAccept_R{}
	res.ErrorCode = func() int32 {
		if len(req.PlrIds) == 0 {
			return Err.Failed
		}

		for i, v := range req.PlrIds {
			ec := plr.GetFriend().Accept(v, req.IsAccept)
			if ec != Err.OK {
				if i == 0 {
					return ec
				} else {
					continue
				}
			}
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
