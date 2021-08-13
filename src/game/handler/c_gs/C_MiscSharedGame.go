package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MiscSharedGame(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MiscSharedGame)

	plr := ctx.(*app.Player)

	res := &msg.GS_MiscSharedGame_R{
		ErrorCode: Err.OK,
	}

	evtmgr.Fire(gconst.Evt_SharedGame, plr, req.Tp)

	plr.SendMsg(res)
}
