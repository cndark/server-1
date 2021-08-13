package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroInherit(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroInherit)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroInherit_R{}
	res.ErrorCode = func() int32 {

		if !gconst.IsHero(req.Id) {
			return Err.Failed
		}

		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		if hero.Id == req.Id {
			return Err.Failed
		}

		conf_h := gamedata.ConfHeroStarUp.Query(hero.Star)
		if conf_h == nil {
			return Err.Failed
		}

		if conf_h.StarInheritNum == 0 {
			return Err.Hero_LowStar
		}

		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return Err.Failed
		}

		conf_m1 := gamedata.ConfMonster.Query(hero.Id)
		if conf_m1 == nil || len(conf_m1.Fragment) == 0 {
			return Err.Failed
		}

		conf_m2 := gamedata.ConfMonster.Query(req.Id)
		if conf_m2 == nil || len(conf_m2.Fragment) == 0 {
			return Err.Failed
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroInherit)
		for _, v := range conf_g.HeroStarInheritCost {
			op.Dec(v.Id, v.N)
		}
		op.Dec(conf_m2.Fragment[0].Id, conf_h.StarInheritNum)
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}
		op.Inc(conf_m1.Fragment[0].Id, conf_h.StarInheritNum)

		res.Rewards = op.Apply().ToMsg()

		hero.Inherit(req.Id)

		return Err.OK
	}()

	plr.SendMsg(res)
}
