package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/act/modules/rushlocal"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActRushLocalGetInfo(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActRushLocalGetInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActRushLocalGetInfo_R{}
	res.ErrorCode = func() int32 {

		res.Data = rushlocal.GetInfo(plr, req.ActName, req.RankId)

		return Err.OK
	}()

	plr.SendMsg(res)
}
