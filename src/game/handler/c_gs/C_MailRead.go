package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MailRead(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MailRead)
	plr := ctx.(*app.Player)

	res := &msg.GS_MailRead_R{}

	func() {
		ec, affected := plr.GetMailBox().Read(req.Id)

		res.ErrorCode = ec

		if ec == Err.OK {
			res.Expire = &msg.MailExpire{
				Id:  affected.Id,
				ETs: affected.ExpireTs.Unix(),
			}
		}
	}()

	plr.SendMsg(res)
}
