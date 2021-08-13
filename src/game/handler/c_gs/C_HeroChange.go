package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

var rand_d = rand.New(rand.NewSource(time.Now().Unix()))

func C_HeroChange(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroChange)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroChange_R{}
	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		if hero.Locked {
			return Err.Hero_Locked
		}

		conf_c := gamedata.ConfHeroChange.Query(hero.Id)
		if conf_c == nil {
			return Err.Failed
		}

		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return Err.Failed
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroChange)
		for _, v := range conf_g.HeroChangeCost {
			if hero.Star == v.Star {
				op.Dec(v.Id, v.N)
			}
		}

		if er := op.CheckEnough(); er != Err.OK {
			return er
		}

		grp := int32(0)
		for _, v := range conf_c.Group {
			if v.Star == hero.Star {
				grp = v.Grp
				break
			}
		}

		id := hero_change_pick(hero.Id, grp)
		conf_h := gamedata.ConfHeroChangeOdds.Query(id)
		if conf_h == nil {
			return Err.Failed
		}

		conf_m := gamedata.ConfMonster.Query(conf_h.Hero)
		if conf_m == nil {
			return Err.Failed
		}

		hero.SetChangeId(conf_h.Hero)

		op.Apply()

		res.ChangeId = conf_h.Hero

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}

func hero_change_pick(heroId int32, grp int32) int32 {
	slt := make(map[int32]int32)

	for _, v := range gamedata.ConfHeroChangeOddsM.Items(grp) {
		if v.Hero != heroId {
			slt[v.Id] += v.Odds
		}
	}

	return utils.PickWeightedMapId(slt)
}
