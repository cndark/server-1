package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MailTakeAttachment(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MailTakeAttachment)
	plr := ctx.(*app.Player)

	res := &msg.GS_MailTakeAttachment_R{}

	func() {
		if ec := plr.GetBag().CheckFull(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		ec, a, affected := plr.GetMailBox().TakeAttachment(req.Id)

		res.ErrorCode = ec

		if ec == Err.OK {
			res.Expire = &msg.MailExpire{
				Id:  affected.Id,
				ETs: affected.ExpireTs.Unix(),
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
