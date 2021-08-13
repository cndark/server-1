package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MailDelOnekey(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_MailDelOnekey)
	plr := ctx.(*app.Player)

	res := &msg.GS_MailDelOnekey_R{}

	res.ErrorCode = Err.OK
	res.Ids = plr.GetMailBox().RemoveOnekey()

	plr.SendMsg(res)
}
