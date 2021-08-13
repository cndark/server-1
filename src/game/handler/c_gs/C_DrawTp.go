package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

func C_DrawTp(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_DrawTp)
	plr := ctx.(*app.Player)

	res := &msg.GS_DrawTp_R{}
	res.ErrorCode = func() int32 {
		conf := gamedata.ConfDraw.Query(req.Tp)
		if conf == nil {
			return Err.Failed
		}

		if !plr.IsModuleOpen(conf.ModuleId) {
			return Err.Plr_ModuleLocked
		}

		// check bag limit
		if er := plr.GetBag().CheckFull(); er != Err.OK {
			return er
		}

		cnt := req.N
		if cnt != int32(1) {
			cnt = int32(10)
		}

		d := plr.GetDraw().GetDrawTp(req.Tp)
		if d == nil {
			return Err.Draw_TpNotFound
		}

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_Draw)

		// 检查免费
		isFree := false
		if cnt == int32(1) && conf.CostFreeTime != 0 &&
			time.Now().After(d.LastTs.Add(time.Duration(conf.CostFreeTime)*time.Second)) {

			if req.Tp == "normal" {
				if plr.GetCounter().GetRemain(gconst.Cnt_DrawNormal_Free) <= 0 {
					isFree = false
				} else {
					cop.DecCounter(gconst.Cnt_DrawNormal_Free, 1)
					isFree = true
					d.LastTs = time.Now()
				}
			} else {
				isFree = true
				d.LastTs = time.Now()
			}
		}

		// 非免费cost
		if !isFree {
			// check cost
			op_check := plr.GetBag().NewOp(gconst.ObjFrom_Draw)
			for _, v := range conf.Cost {
				op_check.Dec(v.Id, v.N*cnt)
			}

			if er1 := op_check.CheckEnough(); er1 != Err.OK {
				if len(conf.EquilCost) == 0 {
					return er1
				}

				discount := int32(10000)
				if cnt == int32(10) { // 十连走折扣
					discount = conf.Discount
				}

				for _, v := range conf.EquilCost {
					cop.Dec(v.Id, int64(float32(v.N*cnt)*float32(discount)/10000.0))
				}

				if er2 := cop.CheckEnough(); er2 != Err.OK {
					return er1
				}
			} else {
				for _, v := range conf.Cost {
					cop.Dec(v.Id, v.N*cnt)
				}
			}
		}

		// draw
		ditem := make(map[int32]int32)
		for i := int32(0); i < cnt; i++ {
			id := plr.GetDraw().DrawOne(req.Tp)

			conf_o := gamedata.ConfDrawOdds.Query(id)
			if conf_o != nil {
				for _, v := range conf_o.Item {
					// 自动分解3星以下英雄
					b := false
					if req.AutoDec && gconst.IsHero(v.Id) {
						conf_m := gamedata.ConfMonster.Query(v.Id)
						if conf_m != nil && conf_m.Star <= 2 {

							conf_s := gamedata.ConfHeroStarUp.Query(conf_m.Star)
							if conf_s != nil {
								for _, s := range conf_s.Sacrifice {
									if s.Type == 1 {
										ditem[s.Id] += s.N * v.N
									}
								}
								b = true
							}
						}
					}

					if !b {
						cop.Inc(v.Id, v.N)
					}

					res.Items = append(res.Items, &msg.Item{Id: v.Id, Num: v.N})
				}
			}
		}

		// 自动分解的需要整合
		for did, dn := range ditem {
			cop.Inc(did, dn)
			res.AutoDecItems = append(res.AutoDecItems, &msg.Item{Id: did, Num: dn})
		}

		// res
		res.Score = plr.GetDraw().Score
		res.DrawTp = d.ToMsg()
		res.Rewards = cop.Apply().ToMsg()

		// fire
		evtmgr.Fire(gconst.Evt_Draw, plr, req.Tp, conf.ModuleId, cnt, res.Items)

		return Err.OK
	}()

	plr.SendMsg(res)
}
