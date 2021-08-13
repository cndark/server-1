package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RelicEat(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RelicEat)
	plr := ctx.(*app.Player)

	res := &msg.GS_RelicEat_R{}

	res.ErrorCode = func() int32 {
		// find master relic
		rlc := plr.GetBag().FindRelic(req.Seq)
		if rlc == nil {
			return Err.Equip_NotFound
		}

		// check full star
		conf_master := gamedata.ConfRelic.Query(rlc.Id)
		if conf_master == nil || rlc.Star >= conf_master.StarLimit {
			return Err.Relic_FullStar
		}

		// check food
		m := make(map[int64]bool)
		for _, seq := range req.EatSeqs {
			if seq == req.Seq {
				return Err.Failed // food is food. not the eater
			}
			m[seq] = true
		}
		if len(m) != len(req.EatSeqs) {
			return Err.Failed
		}

		// eat
		op := plr.GetBag().NewOp(gconst.ObjFrom_RelicEat)
		sum_xp := int32(0)

		for _, seq := range req.EatSeqs {
			food := plr.GetBag().FindRelic(seq)
			if food == nil {
				continue
			}

			conf := gamedata.ConfRelic.Query(food.Id)
			if conf == nil {
				continue
			}

			sum_xp += conf.SupplyExp + food.Star*(food.Star-1)/2*conf.Exp
			op.DelRelic(seq)
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// ok
		op.Apply()

		// add xp
		rlc.AddExp(sum_xp)

		return Err.OK
	}()

	plr.SendMsg(res)
}
