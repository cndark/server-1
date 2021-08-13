package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroDecompose(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroDecompose)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroDecompose_R{}

	res.ErrorCode = func() int32 {
		// limit hero count
		L := len(req.Seqs)
		if L > 30 {
			return Err.Failed
		}

		// check params
		m := make(map[int64]bool)
		for _, seq := range req.Seqs {
			m[seq] = true
		}
		if len(m) != len(req.Seqs) {
			return Err.Failed
		}

		// conf
		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return Err.Failed
		}

		// return cost
		ids := make([]int32, 0, L)

		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroDecompose)

		for _, seq := range req.Seqs {
			// find hero
			hero := plr.GetBag().FindHero(seq)
			if hero == nil || hero.Locked {
				continue
			}

			if plr.GetTeamMgr().InTeam(seq, gconst.TeamType_Dfd) {
				continue
			}

			// 返还 等级消耗
			conf_up := gamedata.ConfHeroUp.Query(hero.Lv)
			if conf_up != nil {
				for _, v := range conf_up.Ret {
					op.Inc(v.Id, float64(v.N)*0.8)
				}
			}

			// 返还星级消耗
			conf_mon := gamedata.ConfMonster.Query(hero.Id)
			if conf_mon != nil {
				cf_len := len(conf_g.CommonFragment)

				// 初始星级对应分解产物
				conf_star := gamedata.ConfHeroStarUp.Query(conf_mon.Star)
				if conf_star != nil {
					for _, v := range conf_star.Sacrifice {
						switch v.Type {
						case 1: // 材料
							op.Inc(v.Id, v.N)
						case 2: // 本体碎片
							op.Inc(conf_mon.Fragment[0].Id, v.N)
						case 3: // 本系通用碎片
							idx := int((conf_mon.Elem-1)*3 + (v.Id - 3))
							if idx >= 0 && idx < cf_len {
								op.Inc(conf_g.CommonFragment[idx], v.N)
							}
						}
					}
				}

				// 升星消耗全额返还
				if conf_mon.Star < hero.Star {
					conf_star2 := gamedata.ConfHeroStarUp.Query(hero.Star)
					if conf_star2 != nil {
						for _, v := range conf_star2.StarRet {
							switch v.Type {
							case 1: // 材料
								op.Inc(v.Id, v.N)
							case 2: // 本体碎片
								op.Inc(conf_mon.Fragment[0].Id, v.N)
							case 3: // 本系通用碎片
								idx := int((conf_mon.Elem-1)*3 + (v.Id - 3))
								if idx >= 0 && idx < cf_len {
									op.Inc(conf_g.CommonFragment[idx], v.N)
								}
							}
						}
					}
				}
			}

			// 返还 饰品消耗
			conf_trk := gamedata.ConfTrinket.Query(hero.Trinket.Lv)
			if conf_trk != nil {
				for _, v := range conf_trk.UpRet {
					op.Inc(v.Id, v.N)
				}
			}

			// del hero
			op.DelHero(seq)

			// added delete seq
			ids = append(ids, hero.Id)
		}

		// apply
		res.Rewards = op.Apply().ToMsg()

		// fire
		evtmgr.Fire(gconst.Evt_HeroDecompose, plr, ids)

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}
