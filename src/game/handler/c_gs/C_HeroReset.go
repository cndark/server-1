package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroReset(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroReset)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroReset_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// check locked
		if hero.Locked {
			return Err.Hero_Locked
		}

		// conf
		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return Err.Failed
		}

		// cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroReset)

		if hero.Lv >= conf_g.HeroResetFreeLv { // 2021-07-02 21级以下免费
			for _, v := range conf_g.HeroResetCost {
				op.Dec(v.Id, v.N)
			}

			if ec := op.CheckEnough(); ec != Err.OK {
				return ec
			}
		}

		// return level cost
		conf_up := gamedata.ConfHeroUp.Query(hero.Lv)
		if conf_up != nil {
			for _, v := range conf_up.Ret {
				op.Inc(v.Id, v.N)
			}
		}

		res.Rewards = op.Apply().ToMsg()

		// reset
		hero.SetLevel(1)
		hero.NewEquipOp().UnequipAll().Apply()

		// fire
		evtmgr.Fire(gconst.Evt_HeroReset, plr, hero.Id)

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}
