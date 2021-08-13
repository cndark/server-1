package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_PlayerHFrameSet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_PlayerHFrameSet)
	plr := ctx.(*app.Player)

	res := &msg.GS_PlayerHFrameSet_R{}
	res.ErrorCode = func() int32 {
		if ec := plr.GetHFrameStore().IsValid(req.Id); ec != Err.OK {
			return ec
		}

		plr.SetHFrame(req.Id)
		res.Id = req.Id

		return Err.OK
	}()

	plr.SendMsg(res)
}
