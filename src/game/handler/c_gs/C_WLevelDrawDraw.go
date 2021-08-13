package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WLevelDrawDraw(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WLevelDrawDraw)
	plr := ctx.(*app.Player)

	res := &msg.GS_WLevelDrawDraw_R{}
	res.ErrorCode = func() int32 {
		ec, items := plr.GetWLevelDraw().Draw(req.Idx)
		if ec != Err.OK {
			return ec
		}

		res.Items = items

		return Err.OK
	}()

	plr.SendMsg(res)
}
