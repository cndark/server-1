package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_CloudGet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_CloudGet)
	plr := ctx.(*app.Player)

	res := &msg.GS_CloudGet_R{}

	res.ErrorCode = func() int32 {
		cld := plr.GetCloud()

		res.Key = req.Key
		res.Val = cld.Get(req.Key)

		return Err.OK
	}()

	plr.SendMsg(res)
}
