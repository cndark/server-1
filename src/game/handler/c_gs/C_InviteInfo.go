package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_InviteInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_InviteInfo)
	plr := ctx.(app.Player)

	res := &msg.GS_InviteInfo_R{}
	res.Info = plr.GetInvite().ToMsg()

	plr.SendMsg(res)
}
