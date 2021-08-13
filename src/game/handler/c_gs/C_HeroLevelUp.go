package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroLevelUp(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroLevelUp)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroLevelUp_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// check N
		if req.N < 1 || req.N > 99999 {
			return Err.Failed
		}

		// check current max level
		to_lv := hero.Lv + req.N
		conf_star := gamedata.ConfHeroStarUp.Query(hero.Star)
		if to_lv > conf_star.MaxLv {
			return Err.Hero_MaxLevelReached
		}

		// check cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroLevelup)

		for i := hero.Lv; i < to_lv; i++ {
			conf_up := gamedata.ConfHeroUp.Query(i)
			if conf_up == nil {
				return Err.Failed
			} else if len(conf_up.Cost) == 0 {
				return Err.Hero_FullLevel
			}

			for _, v := range conf_up.Cost {
				op.Dec(v.Id, v.N)
			}
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		op.Apply()

		// level up
		hero.SetLevel(to_lv)

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}
