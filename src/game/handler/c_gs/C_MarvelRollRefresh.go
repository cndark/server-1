package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MarvelRollRefresh(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MarvelRollRefresh)
	plr := ctx.(*app.Player)

	res := &msg.GS_MarvelRollRefresh_R{}
	res.ErrorCode = func() int32 {

		ec := plr.GetMarvelRoll().Refresh(req.Grp)
		if ec != Err.OK {
			return ec
		}

		group := plr.GetMarvelRoll().Groups[req.Grp]
		if group == nil {
			return Err.MarvelRoll_GroupNotFound
		}

		res.Group = group.ToMsg_Group(req.Grp)

		return Err.OK
	}()

	plr.SendMsg(res)
}
