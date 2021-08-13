package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WLevelOneKeyGJTake(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WLevelOneKeyGJTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_WLevelOneKeyGJTake_R{}
	res.ErrorCode = func() int32 {
		gjLv := plr.GetWLevel().LvNum + 1
		maxLv := gamedata.ConfLimitM.Query().MaxWLevelLv
		if gjLv > maxLv {
			gjLv = maxLv
		}

		conf_w := gamedata.ConfWorldLevelM.Query(gjLv)
		if conf_w == nil {
			return Err.Failed
		}

		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return Err.Failed
		}

		if ec := plr.GetBag().CheckFull(); ec != Err.OK {
			return ec
		}

		n := conf_g.ExploreOnekeyTime * 60

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_WLevelOneKeyGJ)

		cop.DecCounter(gconst.Cnt_WLevelOneKeyGj, 1)
		if ec := cop.CheckEnough(); ec != Err.OK {
			return ec
		}

		for _, v := range conf_w.MinuteCurrency {
			cop.Inc(v.Id, v.N*n)
		}

		for i := int32(0); i < n; i++ {
			for _, v := range utils.Drop(plr, conf_w.ExploreDrop) {
				cop.Inc(v.Id, v.N)
			}
		}

		res.Rewards = cop.Apply().ToMsg()

		evtmgr.Fire(gconst.Evt_WLevelOneKeyGj, plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
