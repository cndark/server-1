package evt

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

func init() {
	// 玩家创号
	evtmgr.On(gconst.Evt_ServerNewUser, func(args ...interface{}) {
		plr := args[1].(*app.Player)

		new_benu_rec(gconst.GLog_UserCreate, plr).
			set("role_name", plr.GetName()).
			set("device_id", plr.GetDevId()).
			set("ip", plr.GetLoginIP()).
			insert()
	})

	// 玩家登陆
	evtmgr.On(gconst.Evt_LoginOnline, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		devid := args[1].(string)
		ip := args[2].(string)

		new_benu_rec(gconst.GLog_UserLogin, plr).
			set("role_name", plr.GetName()).
			set("level", plr.GetLevel()).
			set("ip", ip).
			set("device_id", devid).
			set("device_model", plr.GetModel()).
			set("device_os", plr.GetOs()+plr.GetOsVer()).
			set("game_version", config.Common.Version).
			set("power", plr.GetAtkPwr()).
			insert()
	})

	// 玩家登出
	evtmgr.On(gconst.Evt_LoginOffline, func(args ...interface{}) {
		plr := args[0].(*app.Player)

		dur := plr.GetOfflineTs().Unix() - plr.GetLoginTs().Unix()
		new_benu_rec(gconst.GLog_UserLogoff, plr).
			set("role_name", plr.GetName()).
			set("ip", plr.GetLoginIP()).
			set("online_dur", dur).
			insert()
	})

	// 当前在线人数
	evtmgr.On(gconst.Evt_OnlineNum, func(args ...interface{}) {
		n := args[0].(int32)

		new_benu_rec(gconst.GLog_OnlineNum, nil).
			set("num", n).
			insert()
	})

	// 玩家新手
	evtmgr.On(gconst.Evt_Tutorial, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		tp := args[1].(string)
		step := args[2].(int32)

		new_benu_rec(gconst.GLog_Tutorial, plr).
			set("ip", plr.GetLoginIP()).
			set("step_type", tp).
			set("step_id", step).
			insert()
	})

	// 玩家升级
	evtmgr.On(gconst.Evt_PlrLv, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		lv := args[1].(int32)
		add := args[2].(int32)

		if add > 0 {
			new_benu_rec(gconst.GLog_UserLv, plr).
				set("old_level", lv-add).
				set("new_level", lv).
				insert()
		}
	})

	// 玩家生成订单
	evtmgr.On(gconst.Evt_BillGen, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		pid := args[1].(int32)
		csext := args[2].(string)
		orderid := args[3].(string)

		conf := gamedata.ConfBillProduct.Query(pid)
		if conf == nil {
			return
		}

		new_benu_rec(gconst.GLog_BillGen, plr).
			set("order_id", orderid).
			set("amount", conf.Price).
			set("confid", pid).
			set("payid", conf.PayId).
			set("level", plr.GetLevel()).
			set("vip", plr.GetVipLevel()).
			set("device_id", plr.GetDevId()).
			set("ip", plr.GetLoginIP()).
			set("csext", csext).
			insert()
	})

	// 玩家完成订单
	evtmgr.On(gconst.Evt_BillStats, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		amount := args[1].(int32)
		orderid := args[3].(string)
		pid := args[4].(int32)
		csext := args[5].(string)

		if pid == 0 {
			return
		}

		conf := gamedata.ConfBillProduct.Query(pid)
		if conf == nil {
			return
		}

		new_benu_rec(gconst.GLog_BillFin, plr).
			set("role_name", plr.GetName()).
			set("order_id", orderid).
			set("amount", amount).
			set("confid", pid).
			set("payid", conf.PayId).
			set("level", plr.GetLevel()).
			set("vip", plr.GetVipLevel()).
			set("device_id", plr.GetDevId()).
			set("ip", plr.GetLoginIP()).
			set("csext", csext).
			insert()
	})

	// 玩家背包
	evtmgr.On(gconst.Evt_BagChg, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		chg := args[1].(map[int32]int64)
		from := args[2].(int32)

		rec := new_benu_rec(gconst.GLog_BagChg, plr)
		rec.set("from", from)

		log_chg := []*log_chg_t{}
		for id, n := range chg {
			cur := int64(0)

			tp := gconst.ObjectType(id)
			if tp == gconst.ObjType_Currency {
				cur = n + plr.GetBag().GetCcy(id)
			} else if tp == gconst.ObjType_Item {
				cur = int64(plr.GetBag().GetItem(id))
			}

			log_chg = append(log_chg, &log_chg_t{
				Id:  id,
				Op:  n,
				Cur: cur,
			})
		}
		rec.set("chg", log_chg)

		rec.insert()
	})

	// 玩家抽卡
	evtmgr.On(gconst.Evt_Draw, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		tp := args[1].(string)
		cnt := args[3].(int32)
		items := args[4].([]*msg.Item)

		log_heroes := []*log_draw_hero_t{}
		for _, v := range items {
			if gconst.IsHero(v.Id) {
				conf := gamedata.ConfMonster.Query(v.Id)
				if conf == nil {
					continue
				}

				log_heroes = append(log_heroes, &log_draw_hero_t{
					Id:   v.Id,
					Star: conf.Star,
				})
			}
		}

		new_benu_rec(gconst.GLog_Draw, plr).
			set("draw_tp", tp).
			set("cnt", cnt).
			set("heroes", log_heroes).
			insert()
	})

	// 玩家商场购买
	evtmgr.On(gconst.Evt_ShopBuy, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		shopid := args[1].(int32)
		shopitem := args[2].(map[int32]int32)

		for id, n := range shopitem {
			conf := gamedata.ConfShopItem.Query(id)
			if conf == nil {
				continue
			}

			new_benu_rec(gconst.GLog_ShopBuy, plr).
				set("shop_id", core.I32toa(shopid)).
				set("goods_id", core.I32toa(conf.Item)).
				set("goods_num", conf.Num*n).
				set("currency", conf.Currency).
				set("discount", conf.Discount/1000).
				insert()
		}
	})

	// 玩家聊天
	evtmgr.On(gconst.Evt_SendChat, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		tp := args[1].(int32)
		content := args[2].(string)
		toid := args[3].(string)
		toName := args[4].(string)

		new_benu_rec(gconst.GLog_Chat, plr).
			set("role_name", plr.GetName()).
			set("chat_channel_id", tp).
			set("chat_role_id", toid).
			set("chat_role_name", toName).
			set("chat_msg", content).
			set("chat_time", time.Now().Unix()).
			set("guild_id", plr.GetGuildId()).
			set("guild_name", plr.GetGuildName()).
			insert()
	})

	// 玩家推图通过
	evtmgr.On(gconst.Evt_WLevelLv, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		lvNum := args[1].(int32)
		addNum := args[2].(int32)

		if addNum > 0 {
			new_benu_rec(gconst.GLog_WLevel, plr).
				set("lv_num", lvNum).
				insert()
		}
	})

	// 玩家推图战斗
	evtmgr.On(gconst.Evt_WLevelFight, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		lvNum := args[1].(int32)
		winner := args[2].(int32)
		team := args[3].(*msg.TeamFormation)

		if winner == 1 {
			return
		}

		log_team := team_log(plr, team)
		new_benu_rec(gconst.GLog_WLevelFailed, plr).
			set("lv_num", lvNum).
			set("team", log_team).
			insert()
	})

	// 玩家爬塔
	evtmgr.On(gconst.Evt_TowerLv, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		lvNum := args[1].(int32)

		new_benu_rec(gconst.GLog_TowerLv, plr).
			set("lv_num", lvNum).
			insert()
	})

	// 玩家爬塔战斗
	evtmgr.On(gconst.Evt_TowerFight, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		lvNum := args[1].(int32)
		winner := args[2].(int32)
		team := args[3].(*msg.TeamFormation)

		if winner == 1 {
			return
		}

		log_team := team_log(plr, team)
		new_benu_rec(gconst.GLog_TowerFailed, plr).
			set("lv_num", lvNum).
			set("team", log_team).
			insert()
	})

	// 玩家爬塔扫荡
	evtmgr.On(gconst.Evt_TowerRaid, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		lvNum := args[1].(int32)

		new_benu_rec(gconst.GLog_TowerRaid, plr).
			set("lv_num", lvNum).
			insert()
	})

	// 玩家积分竞技场打架
	evtmgr.On(gconst.Evt_ArenaFight, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		isWin := args[1].(bool)
		isBot := args[2].(bool)
		subAtk := args[3].(int32)
		score := args[4].(int32)

		new_benu_rec(gconst.GLog_ArenaFight, plr).
			set("is_win", isWin).
			set("is_bot", isBot).
			set("sub_atk", subAtk).
			set("add_score", score).
			insert()
	})

}
