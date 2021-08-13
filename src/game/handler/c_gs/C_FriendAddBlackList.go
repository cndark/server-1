package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_FriendAddBlackList(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_FriendAddBlackList)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendAddBlackList_R{}
	res.ErrorCode = func() int32 {
		if req.IsAdd {
			ec := plr.GetFriend().AddBlackList(req.PlrId)
			if ec != Err.OK {
				return ec
			}
		} else {
			ec := plr.GetFriend().RemoveBlackList(req.PlrId)
			if ec != Err.OK {
				return ec
			}
		}

		res.IsAdd = req.IsAdd
		res.PlrId = req.PlrId

		return Err.OK
	}()

	plr.SendMsg(res)
}
