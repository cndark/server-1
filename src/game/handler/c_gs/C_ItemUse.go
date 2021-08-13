package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/itemuse"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ItemUse(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ItemUse)
	plr := ctx.(*app.Player)

	res := &msg.GS_ItemUse_R{ErrorCode: Err.OK}

	f := itemuse.GetUseHandler(req.Id)

	if f == nil || req.N <= 0 || req.N > 99999 {
		res.ErrorCode = Err.Failed
	}

	if res.ErrorCode == Err.OK {
		res = f(plr, req.Id, req.N)
	}

	plr.SendMsg(res)
}
