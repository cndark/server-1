package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArmorCompose(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArmorCompose)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArmorCompose_R{}

	res.ErrorCode = func() int32 {
		bag := plr.GetBag()

		// check
		conf := gamedata.ConfItem.Query(req.SrcId)
		if conf == nil || len(conf.ArmorCompose) == 0 {
			return Err.Armor_Uncomposable
		}

		if req.N <= 0 {
			return Err.Failed
		}

		// compose
		op := bag.NewOp(gconst.ObjFrom_ArmorCompose)

		// cost
		op.Dec(req.SrcId, req.N*3)
		op.Inc(req.SrcId+1, req.N)

		for _, v := range conf.ArmorCompose {
			op.Dec(v.Id, v.N)
		}

		// check
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// apply
		res.Rewards = op.Apply().ToMsg()

		// fire
		evtmgr.Fire(gconst.Evt_ArmorCompose, plr, int32(1))

		return Err.OK
	}()

	plr.SendMsg(res)
}
