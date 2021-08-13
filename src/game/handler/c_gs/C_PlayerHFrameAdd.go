package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_PlayerHFrameAdd(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_PlayerHFrameAdd)
	plr := ctx.(*app.Player)

	res := &msg.GS_PlayerHFrameAdd_R{}
	res.ErrorCode = func() int32 {
		conf := gamedata.ConfPlayerHFrame.Query(req.Id)
		if conf == nil {
			return Err.Failed
		}

		if ec := utils.IsStatusTabArrayConform(plr, conf.Status); ec != Err.OK {
			return ec
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_PlrChangeHFrame)
		for _, v := range conf.Cost {
			op.Dec(v.Id, v.N)
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		if ec := plr.GetHFrameStore().Add(req.Id); ec != Err.OK {
			return ec
		}

		op.Apply()

		plr.SetHFrame(req.Id)
		res.Id = req.Id

		return Err.OK
	}()

	plr.SendMsg(res)
}
