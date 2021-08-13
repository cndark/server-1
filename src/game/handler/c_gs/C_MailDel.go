package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_MailDel(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MailDel)
	plr := ctx.(*app.Player)

	res := &msg.GS_MailDel_R{}

	func() {
		ec := plr.GetMailBox().Remove(req.Id)

		res.ErrorCode = ec
		res.Id = req.Id
	}()

	plr.SendMsg(res)
}
