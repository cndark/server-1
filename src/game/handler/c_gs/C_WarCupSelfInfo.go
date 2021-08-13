package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupSelfInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WarCupSelfInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupSelfInfo_R{}
	res.ErrorCode = func() int32 {

		res.VsData, res.CurReplay = warcup.WarCupSelfVsInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
