package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_PlayerHeadSet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_PlayerHeadSet)
	plr := ctx.(*app.Player)

	res := &msg.GS_PlayerHeadSet_R{}
	res.ErrorCode = func() int32 {
		L := len(req.Head)
		if L == 0 || L > 200 {
			return Err.Failed
		}

		plr.User().Head = req.Head

		return Err.OK
	}()

	plr.SendMsg(res)
}
