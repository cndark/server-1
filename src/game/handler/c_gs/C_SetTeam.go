package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_SetTeam(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_SetTeam)
	plr := ctx.(*app.Player)

	res := &msg.GS_SetTeam_R{}
	res.ErrorCode = func() int32 {
		if !plr.IsTeamFormationValid(req.T) {
			return Err.Plr_TeamInvalid
		}

		plr.GetTeamMgr().SetTeam(req.Tp, req.T)

		res.Tp = req.Tp
		res.T = req.T

		return Err.OK
	}()

	plr.SendMsg(res)
}
