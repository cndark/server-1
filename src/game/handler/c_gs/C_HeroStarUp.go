package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroStarUp(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroStarUp)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroStarUp_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// check full star
		conf_mon := gamedata.ConfMonster.Query(hero.Id)
		if hero.Star+1 > conf_mon.StarLimit {
			return Err.Hero_FullStar
		}

		// get conf
		conf_star := gamedata.ConfHeroStarUp.Query(hero.Star)
		if conf_star == nil || len(conf_star.StarCost) == 0 {
			return Err.Hero_FullStar
		}

		// ----------- CHEAT -----------
		// pt-bot 直接升星
		if plr.GetSdk() == "soda.pressure" {
			hero.SetStar(hero.Star + 1)
			return Err.OK
		}
		// ----------- CHEAT -----------

		// check param
		if len(req.Cost) != len(conf_star.StarCost) {
			return Err.Failed
		}

		// check elem
		elem_me := hero.GetElem()
		if elem_me == 0 {
			return Err.Failed
		}

		// check recipe
		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroStar)

		type info_t struct {
			hero_lv int32
			trk_lv  int32
		}
		m := make(map[int64]*info_t)

		for i, v := range conf_star.StarCost {
			harr := req.Cost[i]

			// check count
			if len(harr.Seqs) != int(v.N) {
				return Err.Failed
			}

			// check type, star
			for _, seq := range harr.Seqs {
				h := plr.GetBag().FindHero(seq)
				if h == nil {
					return Err.Hero_NotFound
				}

				if plr.GetTeamMgr().InTeam(seq, gconst.TeamType_Dfd) {
					return Err.Plr_InTeam
				}

				switch v.Tp {
				case 1: // 自己
					if m[seq] != nil || h.Id != hero.Id || h.Star != v.Star {
						return Err.Hero_InvalidCostHero
					}

					m[seq] = &info_t{h.Lv, h.Trinket.Lv}
					op.DelHero(seq)

				case 3: // 同系英雄
					if m[seq] != nil || h.GetElem() != elem_me || h.Star != v.Star {
						return Err.Hero_InvalidCostHero
					}

					m[seq] = &info_t{h.Lv, h.Trinket.Lv}
					op.DelHero(seq)

				case 4: // 任意英雄
					if m[seq] != nil || h.Star != v.Star {
						return Err.Hero_InvalidCostHero
					}

					m[seq] = &info_t{h.Lv, h.Trinket.Lv}
					op.DelHero(seq)
				}
			}
		}

		// check addcost
		for _, v := range conf_star.StarAddCost {
			op.Dec(v.Id, v.N)
		}

		// check
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// 返还 <消耗英雄> 的等级材料
		for _, info := range m {
			// 返还 <消耗英雄> 的等级材料
			conf_heroup := gamedata.ConfHeroUp.Query(info.hero_lv)
			if conf_heroup != nil {
				for _, r := range conf_heroup.Ret {
					op.Inc(r.Id, r.N)
				}
			}

			// 返还 <消耗英雄> 的饰品升级材料
			conf_trk := gamedata.ConfTrinket.Query(info.trk_lv)
			if conf_trk != nil {
				for _, r := range conf_trk.UpRet {
					op.Inc(r.Id, r.N)
				}
			}
		}

		// apply
		res.Rewards = op.Apply().ToMsg()

		// star up
		hero.SetStar(hero.Star + 1)

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}
