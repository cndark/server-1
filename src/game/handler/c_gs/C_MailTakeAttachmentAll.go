package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MailTakeAttachmentAll(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_MailTakeAttachmentAll)
	plr := ctx.(*app.Player)

	res := &msg.GS_MailTakeAttachmentAll_R{}

	func() {
		if ec := plr.GetBag().CheckFull(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		ec, a, affected := plr.GetMailBox().TakeAttachmentAll()

		res.ErrorCode = ec

		if ec == Err.OK {
			for _, m := range affected {
				res.Expires = append(res.Expires, &msg.MailExpire{
					Id:  m.Id,
					ETs: m.ExpireTs.Unix(),
				})
			}

			// add attachment
			op := plr.GetBag().NewOp(gconst.ObjFrom_Mail)
			for _, res := range a {
				op.Inc(res.Id, res.N)
			}
			op.Apply()
		}
	}()

	plr.SendMsg(res)
}
