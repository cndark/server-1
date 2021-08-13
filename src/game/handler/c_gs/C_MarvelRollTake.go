package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MarvelRollTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MarvelRollTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_MarvelRollTake_R{}
	res.ErrorCode = func() int32 {

		ec, ids := plr.GetMarvelRoll().Take(req.Grp, req.IsTen)
		if ec != Err.OK {
			return ec
		}

		res.Grp = req.Grp
		res.Ids = ids

		return Err.OK
	}()

	plr.SendMsg(res)
}
