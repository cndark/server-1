package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ItemChoose(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ItemChoose)
	plr := ctx.(*app.Player)

	res := &msg.GS_ItemChoose_R{}
	res.ErrorCode = func() int32 {
		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemChoose)

		conf := gamedata.ConfItemChoose.Query(req.Id)
		if conf == nil {
			return Err.Failed
		}

		for idx, n := range req.Idx {
			if n <= 0 || n > 999 {
				n = 1
			}

			if idx <= 0 || idx > int32(len(conf.Reward)) {
				return Err.Failed
			}

			op.Dec(req.Id, int64(n))

			op.Inc(conf.Reward[idx-1].Id, int64(conf.Reward[idx-1].N*int64(n)))
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		res.Rewards = op.Apply().ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}
