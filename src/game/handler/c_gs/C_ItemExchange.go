package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ItemExchange(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ItemExchange)
	plr := ctx.(*app.Player)

	res := &msg.GS_ItemExchange_R{}
	res.ErrorCode = func() int32 {
		if req.N <= 0 || req.N > 9999 {
			req.N = 1
		}

		conf := gamedata.ConfItemExchange.Query(req.Id)
		if conf == nil {
			return Err.Failed
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemExchange)
		for _, v := range conf.Cost {
			op.Dec(v.Id, v.N*int64(req.N))
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		op.Inc(req.Id, req.N)

		res.Rewards = op.Apply().ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}
