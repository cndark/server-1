package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupAuditionRank(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WarCupAuditionRank)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupAuditionRank_R{}
	res.ErrorCode = func() int32 {

		res.Rows = warcup.WarCupAuditonRank(req.Top, req.N)

		return Err.OK
	}()

	plr.SendMsg(res)
}
