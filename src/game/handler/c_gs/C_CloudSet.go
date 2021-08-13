package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_CloudSet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_CloudSet)
	plr := ctx.(*app.Player)

	res := &msg.GS_CloudSet_R{}

	res.ErrorCode = func() int32 {
		cld := plr.GetCloud()

		if !cld.Set(req.Key, req.Val) {
			return Err.Cld_InvalidKey
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
