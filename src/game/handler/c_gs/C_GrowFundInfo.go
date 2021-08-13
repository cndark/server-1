package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/growfund"
	"fw/src/game/msg"
)

func C_GrowFundInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GrowFundInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_GrowFundInfo_R{
		SvrBuyCnt: growfund.GrowFundSvr.SvrCnt,
	}

	plr.SendMsg(res)
}
