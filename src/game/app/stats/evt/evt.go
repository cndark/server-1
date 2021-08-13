package evt

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/sched/async"
	"fw/src/game/app"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

func flush(f func()) {
	async.PushQ(gconst.AQ_Stats, f)
}

// ============================================================================

func init() {
	// 登录
	//	area 区域 Id
	//	svr  服务器名称
	//  sdk  sdk
	//	uid  玩家 Id
	//	cts  玩家创建时间
	//	lts  当天最后一次登录时间
	//	day  第几天
	//	n    当天登录次数
	evtmgr.On(gconst.Evt_LoginOnline, func(args ...interface{}) {
		plr := args[0].(*app.Player)

		area := config.Common.Area.Id
		svr := config.CurGame.Name
		sdk := plr.GetSdk()
		uid := plr.GetId()
		cts := plr.GetCreateTs()
		lts := plr.GetLoginTs()

		cdate := core.StartOfDay(cts)
		day := int32(time.Now().Sub(cdate).Hours())/24 + 1

		flush(func() {
			dbmgr.DBStats.UpsertByCond(
				dbmgr.C_tabname_login,
				db.M{
					"uid": uid,
					"day": day,
				},
				db.M{
					"$setOnInsert": db.M{
						"area": area,
						"sdk":  sdk,
						"cts":  cts,
					},
					"$set": db.M{
						"svr": svr,
						"lts": lts,
					},
					"$inc": db.M{
						"n": 1,
					},
				},
			)
		})
	})

	// 充值
	//	area 区域 Id
	//  svr  服务器名称
	//	sdk  sdk
	//  uid  玩家 Id
	//  name 玩家最新名字
	//  cts  玩家创建时间
	//  bts  当天最后一次充值时间
	//  day  第几天
	//  amt  当天充值总额
	//  ccy  币种
	//  n    当天充值次数
	evtmgr.On(gconst.Evt_BillStats, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		amt := args[1].(int32)
		ccy := args[2].(string)

		area := config.Common.Area.Id
		svr := config.CurGame.Name
		sdk := plr.GetSdk()
		uid := plr.GetId()
		name := plr.GetName()
		cts := plr.GetCreateTs()

		cdate := core.StartOfDay(cts)
		now := time.Now()

		day := int32(now.Sub(cdate).Hours())/24 + 1

		n := 1
		if amt == 0 {
			n = 0
		}

		flush(func() {
			dbmgr.DBStats.UpsertByCond(
				dbmgr.C_tabname_bill,
				db.M{
					"uid": uid,
					"day": day,
				},
				db.M{
					"$setOnInsert": db.M{
						"area": area,
						"sdk":  sdk,
						"cts":  cts,
					},
					"$set": db.M{
						"svr":  svr,
						"name": name,
						"bts":  now,
						"ccy":  ccy,
					},
					"$inc": db.M{
						"amt": amt,
						"n":   n,
					},
				},
			)
		})
	})

	// 每天在线人数
	evtmgr.On(gconst.Evt_OnlineNum, func(args ...interface{}) {
		avg := args[1].(int32)  // 平均在线人数
		peek := args[2].(int32) // 最高在线人数

		// flush
		area := config.Common.Area.Id
		svr := config.CurGame.Name
		ts := core.StartOfDay(time.Now())

		flush(func() {
			dbmgr.DBStats.UpsertByCond(
				dbmgr.C_tabname_online,
				db.M{
					"area": area,
					"svr":  svr,
					"ts":   ts,
				},
				db.M{
					"$set": db.M{
						"avg":  avg,
						"peek": peek,
					},
				},
			)
		})
	})

	// 引导
	//	area  区域 Id
	//  svr   服务器名称
	//	sdk   sdk
	//  uid   玩家 Id
	//  cts   玩家创建时间
	//  ots   最后操作时间
	//  tut   类型步骤
	evtmgr.On(gconst.Evt_Tutorial, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		tp := args[1].(string)
		step := args[2].(int32)

		area := config.Common.Area.Id
		svr := config.CurGame.Name
		sdk := plr.GetSdk()
		uid := plr.GetId()
		cts := plr.GetCreateTs()
		ots := time.Now()
		tut := "tut." + tp

		flush(func() {
			dbmgr.DBStats.UpsertByCond(
				dbmgr.C_tabname_tutorial,
				db.M{
					"uid": uid,
				},
				db.M{
					"$setOnInsert": db.M{
						"area": area,
						"sdk":  sdk,
						"cts":  cts,
					},
					"$set": db.M{
						"svr": svr,
						"ots": ots,
						tut:   step,
					},
				},
			)
		})
	})

	// 关卡
	//	area  区域 Id
	//  svr   服务器名称
	//	sdk   sdk
	//  uid   玩家 Id
	//  cts   玩家创建时间
	//  ots   最后操作时间
	//  wlv   当前关卡lv
	evtmgr.On(gconst.Evt_WLevelLv, func(args ...interface{}) {
		plr := args[0].(*app.Player)
		wlv := args[1].(int32)
		addNum := args[2].(int32)

		if addNum <= 0 {
			return
		}

		area := config.Common.Area.Id
		svr := config.CurGame.Name
		sdk := plr.GetSdk()
		uid := plr.GetId()
		cts := plr.GetCreateTs()
		ots := time.Now()

		flush(func() {
			dbmgr.DBStats.UpsertByCond(
				dbmgr.C_tabname_wlevel,
				db.M{
					"uid": uid,
				},
				db.M{
					"$setOnInsert": db.M{
						"area": area,
						"sdk":  sdk,
						"cts":  cts,
					},
					"$set": db.M{
						"svr": svr,
						"ots": ots,
						"wlv": wlv,
					},
				},
			)
		})
	})

}
