package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_DrawGetInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_DrawGetInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_DrawGetInfo_R{}
	res.ErrorCode = func() int32 {
		res.Draw = plr.GetDraw().ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}
