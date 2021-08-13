package cond

import (
	"fw/src/game/app/comp"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
)

// ============================================================================
// 条件id, 策划定好id后手动填写, 包含玩家和家族的条件

const (
	// 玩家
	plrcond_plr_login        = 1  // 登录天数
	plrcond_cost_ccy_cnt     = 2  // 累计消费货币
	plrcond_friend_gift      = 3  // 好友赠送体力
	plrcond_appoint_fin      = 4  // 派遣完成
	plrcond_shopbuy          = 5  // 商店购买次数
	plrcond_shared_game      = 6  // 分享游戏
	plrcond_arrmor_compose   = 7  // 装备合成
	plrcond_wlevel_gj        = 8  // 领取挂机
	plrcond_marvelroll       = 9  // 奇迹之盘抽奖
	plrcond_add_hero         = 10 // 获得英雄数量
	plrcond_draw             = 11 // 召唤
	plrcond_goldhand_cnt     = 12 // 点金次数
	plrcond_taskdaily_fin    = 13 // 日常任务完成
	plrcond_bill_baseccy_cnt = 14 // 充值n基准货币
	plrcond_add_friend       = 15 // 添加好友
	plrcond_pushgift         = 16 // 推送礼包购买
	plrcond_ranklike         = 17 // 排行榜点赞
	plrcond_add_hero_draw    = 18 // 抽卡召唤获得英雄数量
	plrcond_add_relic        = 19 // 获得神器数量

	plrcond_cost_item_cnt = 20 // 累计物品

	plrcond_targetdays_cnt = 30 // 开服庆典任务领取数
	plrcond_taskmonth_cnt  = 31 // 每月任务完成数量
	plrcond_appoint_send   = 44 // 派遣次数

	plrcond_plr_lv              = 101 // 玩家等级
	plrcond_plr_atkpwr          = 102 // 玩家历史最大战力
	plrcond_hero_reach_star_cnt = 103 // 英雄星级p1个数
	plrcond_plr_viplv           = 104 // 玩家vip等级
	plrcond_decompose_hero      = 105 // 分解英雄个数
	plrcond_hero_reach_lv_cnt   = 106 // 英雄等级p1个数
	plrcond_hero_reset_cnt      = 107 // 英雄重生次数
	plrcond_trinket_lv_cnt      = 108 // 升级饰品次数

	plrcond_wlevel_lv              = 201 // 关卡到达
	plrcond_wlevel_fight_cnt       = 202 // 关卡战斗次数
	plrcond_arena_rank             = 203 // 竞技场排名
	plrcond_arena_max_score        = 204 // 竞技场最高积分
	plrcond_tower_lv               = 205 // 爬塔到达
	plrcond_tower_fight_cnt        = 206 // 爬塔战斗次数
	plrcond_arena_fight_cnt        = 207 // 竞技场挑战次数
	plrcond_guildboss_fight_cnt    = 208 // 家族boss战斗次数
	plrcond_crusade_fight_cnt      = 209 // 英灵试炼战斗次数
	plrcond_onekeygj_cnt           = 210 // 快速挂机次数
	plrcond_tower_raid_cnt         = 211 // 爬塔扫荡次数
	plrcond_rift_monster_fight_cnt = 212 // 击败裂隙怪物次数
	plrcond_rift_mine_occupy_cnt   = 213 // 击败裂隙矿次数
	plrcond_rift_mine_take_cnt     = 214 // 收获裂隙矿产奖励
	plrcond_rift_box_take_cnt      = 215 // 收获裂隙宝箱奖励
	plrcond_ladder_fight_cnt       = 220 // 烈焰天梯战斗次数

	plrcond_guild_donate_cnt = 301 // 家族捐献次数
	plrcond_guild_wish_cnt   = 302 // 家族许愿次数
	plrcond_guild_help_cnt   = 303 // 家族碎片助力次数
	plrcond_guild_fin_order  = 304 // 家族完成港口订单
	plrcond_guild_join_cnt   = 306 // 家族申请次数
	plrcond_guild_sign_cnt   = 307 // 家族签到次数

	plrcond_warcup_guess_cnt     = 401 // 本服杯赛竞猜次数
	plrcond_warcup_guess_win_cnt = 402 // 本服杯赛竞猜正确次数
	plrcond_warcup_chat_cnt      = 403 // 本服杯赛聊天发言次数
	plrcond_warcup_watch_cnt     = 404 // 本服杯赛观看比赛次数

	plrcond_actmaze_cnt       = 411 // 迷宫点击次数
	plrcond_actmaze_event_cnt = 412 // 迷宫事件次数
	plrcond_actmaze_score     = 413 // 迷宫积分

)

// ============================================================================

// 监听回调函数
type cond_cb_t struct {
	Evt string
	F   func(p1 int32, cv ICondObj, args []interface{})
}

var cond_cb = map[int32]*cond_cb_t{

	// 玩家: 玩家登录
	plrcond_plr_login: {

		gconst.Evt_PlrDailyOnline,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家: 玩家等级
	plrcond_plr_lv: {

		gconst.Evt_PlrLv,

		func(p1 int32, cv ICondObj, args []interface{}) {
			lv := float64(args[1].(int32))

			if cv.GetVal() < lv {
				cv.SetVal(lv)
			}

			return
		},
	},

	// 玩家: 玩家历史最大战力
	plrcond_plr_atkpwr: {

		gconst.Evt_PlrAtkPwr,

		func(p1 int32, cv ICondObj, args []interface{}) {
			atkpwr := float64(args[1].(int32))

			if atkpwr > cv.GetVal() {
				cv.SetVal(atkpwr)
			}

			return
		},
	},

	// 玩家: 累计消费货币
	plrcond_cost_ccy_cnt: {

		gconst.Evt_CcyDel,

		func(p1 int32, cv ICondObj, args []interface{}) {
			cost := args[1].(map[int32]int64)

			n := float64(cost[p1])
			if n >= 0 {
				return
			}

			cv.AddVal(-n)

			return
		},
	},

	// 玩家: 累计消费道具
	plrcond_cost_item_cnt: {

		gconst.Evt_ItemDel,

		func(p1 int32, cv ICondObj, args []interface{}) {
			cost := args[1].(map[int32]int32)

			n := float64(cost[p1])
			if n >= 0 {
				return
			}

			cv.AddVal(-n)

			return
		},
	},

	// 玩家: 好友赠送体力
	plrcond_friend_gift: {

		gconst.Evt_FriendGive,

		func(p1 int32, cv ICondObj, args []interface{}) {
			n := args[1].(int32)

			cv.AddVal(float64(n))

			return
		},
	},

	// 玩家: 派遣完成
	plrcond_appoint_fin: {

		gconst.Evt_AppointFin,

		func(p1 int32, cv ICondObj, args []interface{}) {
			star := args[1].(int32)

			if p1 == -1 || star == p1 {
				cv.AddVal(1)
			}

			return
		},
	},

	// 玩家: 商店购买
	plrcond_shopbuy: {

		gconst.Evt_ShopBuy,

		func(p1 int32, cv ICondObj, args []interface{}) {
			shopid := args[1].(int32)

			if p1 == -1 || p1 == shopid {
				cv.AddVal(1)
			}

			return
		},
	},

	plrcond_shared_game: {

		gconst.Evt_SharedGame,

		func(p1 int32, cv ICondObj, args []interface{}) {
			tp := args[1].(int32)

			if p1 == -1 || tp == p1 {
				cv.AddVal(1)
			}

			return
		},
	},

	// 玩家: 装备合成
	plrcond_arrmor_compose: {

		gconst.Evt_ArmorCompose,

		func(p1 int32, cv ICondObj, args []interface{}) {
			n := args[1].(int32)

			cv.AddVal(float64(n))

			return
		},
	},

	// 玩家: 领取挂机
	plrcond_wlevel_gj: {

		gconst.Evt_WLevelGj,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家: 奇迹之盘抽奖次数
	plrcond_marvelroll: {

		gconst.Evt_MarvelRoll,

		func(p1 int32, cv ICondObj, args []interface{}) {
			mid := args[2].(int32)
			n := args[3].(int32)

			if p1 == -1 || mid == p1 {
				cv.AddVal(float64(n))
			}

			return
		},
	},

	// 玩家: 获得英雄
	plrcond_add_hero: {

		gconst.Evt_HeroAdd,

		func(p1 int32, cv ICondObj, args []interface{}) {
			heroes := args[1].((map[int64]*comp.Hero))

			if p1 == -1 {
				cv.AddVal(float64(len(heroes)))
			} else {
				n := 0
				for _, hero := range heroes {
					conf := gamedata.ConfMonster.Query(hero.Id)
					if conf != nil && conf.Star == p1 {
						n++
					}
				}

				if n > 0 {
					cv.AddVal(float64(n))
				}
			}

			return
		},
	},

	// 玩家: 抽卡召唤获得英雄
	plrcond_add_hero_draw: {

		gconst.Evt_HeroAdd,

		func(p1 int32, cv ICondObj, args []interface{}) {
			heroes := args[1].((map[int64]*comp.Hero))
			from := args[2].(int32)

			if from != gconst.ObjFrom_Draw {
				return
			}

			if p1 == -1 {
				cv.AddVal(float64(len(heroes)))
			} else {
				n := 0
				for _, hero := range heroes {
					conf := gamedata.ConfMonster.Query(hero.Id)
					if conf != nil && conf.Star == p1 {
						n++
					}
				}

				if n > 0 {
					cv.AddVal(float64(n))
				}
			}

			return
		},
	},

	// 玩家: 获得神器
	plrcond_add_relic: {

		gconst.Evt_RelicAdd,

		func(p1 int32, cv ICondObj, args []interface{}) {
			relics := args[1].((map[int64]*comp.Relic))

			if p1 == -1 {
				cv.AddVal(float64(len(relics)))
			} else {
				n := 0
				for _, rc := range relics {
					conf := gamedata.ConfRelic.Query(rc.Id)
					if conf != nil && conf.Color == p1 {
						n++
					}
				}

				if n > 0 {
					cv.AddVal(float64(n))
				}
			}

			return
		},
	},

	// 玩家: 奇迹之盘抽奖次数
	plrcond_draw: {

		gconst.Evt_Draw,

		func(p1 int32, cv ICondObj, args []interface{}) {
			mid := args[2].(int32)
			n := args[3].(int32)

			if p1 == -1 &&
				(mid == gconst.ModuleId_DrawNormal || mid == gconst.ModuleId_DrawHigh) {
				cv.AddVal(float64(n))
			} else if p1 == mid {
				cv.AddVal(float64(n))
			}

			return
		},
	},

	// 玩家: 点金手次数
	plrcond_goldhand_cnt: {

		gconst.Evt_GoldHand,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家: 日常任务完成
	plrcond_taskdaily_fin: {

		gconst.Evt_TaskDaily_Fin,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家: 充值基准货币数
	plrcond_bill_baseccy_cnt: {

		gconst.Evt_BillDone,

		func(p1 int32, cv ICondObj, args []interface{}) {
			pid := args[1].(int32)

			prod := gamedata.ConfBillProduct.Query(pid)
			if prod == nil {
				return
			}

			cv.AddVal(float64(prod.BaseCcy))

			return
		},
	},

	// 玩家: 添加好友
	plrcond_add_friend: {

		gconst.Evt_FriendAdd,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家: 推送礼包购买
	plrcond_pushgift: {

		gconst.Evt_PushGift,

		func(p1 int32, cv ICondObj, args []interface{}) {
			id := args[1].(int32)

			if p1 == -1 || id == p1 {
				cv.AddVal(1)
			}

			return
		},
	},

	// 玩家: 排行榜点赞
	plrcond_ranklike: {

		gconst.Evt_RankLike,

		func(p1 int32, cv ICondObj, args []interface{}) {
			rkid := args[1].(int32)

			if p1 == -1 || rkid == p1 {
				cv.AddVal(1)
			}

			return
		},
	},

	// 玩家: 英雄星级p1个数
	plrcond_hero_reach_star_cnt: {

		gconst.Evt_HeroStar,

		func(p1 int32, cv ICondObj, args []interface{}) {
			plr := args[0].(IPlayer)

			n := float64(plr.GetBag().HeroReachStarCnt(p1))

			if cv.GetVal() < n {
				cv.SetVal(n)
			}

			return
		},
	},

	// 玩家: 玩家vip等级
	plrcond_plr_viplv: {

		gconst.Evt_VipLv,

		func(p1 int32, cv ICondObj, args []interface{}) {
			lv := float64(args[1].(int32))

			if cv.GetVal() < lv {
				cv.SetVal(lv)
			}

			return
		},
	},

	// 玩家: 分解英雄个数
	plrcond_decompose_hero: {

		gconst.Evt_HeroDecompose,

		func(p1 int32, cv ICondObj, args []interface{}) {
			ids := args[1].([]int32)

			n := len(ids)
			if n > 0 {
				cv.AddVal(float64(n))
			}

			return
		},
	},

	// 玩家: 英雄等级p1个数
	plrcond_hero_reach_lv_cnt: {

		gconst.Evt_HeroLv,

		func(p1 int32, cv ICondObj, args []interface{}) {
			plr := args[0].(IPlayer)

			n := float64(plr.GetBag().HeroReachLvCnt(p1))

			if cv.GetVal() < n {
				cv.SetVal(n)
			}

			return
		},
	},

	// 玩家: 英雄重生次数
	plrcond_hero_reset_cnt: {

		gconst.Evt_HeroReset,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家: 升级饰品次数
	plrcond_trinket_lv_cnt: {

		gconst.Evt_HeroTrinketLv,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家：关卡到达
	plrcond_wlevel_lv: {

		gconst.Evt_WLevelLv,

		func(p1 int32, cv ICondObj, args []interface{}) {
			lvNum := args[1].(int32)

			if cv.GetVal() < float64(lvNum) {
				cv.SetVal(float64(lvNum))
			}

			return
		},
	},

	// 玩家：增加关卡战斗次数
	plrcond_wlevel_fight_cnt: {

		gconst.Evt_WLevelFight,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家：爬塔到达
	plrcond_tower_lv: {

		gconst.Evt_TowerLv,

		func(p1 int32, cv ICondObj, args []interface{}) {
			lvNum := args[1].(int32)

			if cv.GetVal() < float64(lvNum) {
				cv.SetVal(float64(lvNum))
			}

			return
		},
	},

	// 玩家：爬塔次数战斗次数
	plrcond_tower_fight_cnt: {

		gconst.Evt_TowerFight,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家：竞技场最高积分
	plrcond_arena_max_score: {

		gconst.Evt_ArenaScore,

		func(p1 int32, cv ICondObj, args []interface{}) {

			score := float64(args[1].(int32))

			if cv.GetVal() < score {
				cv.SetVal(score)
			}

			return
		},
	},

	// 玩家：竞技场战斗次数
	plrcond_arena_fight_cnt: {

		gconst.Evt_ArenaFight,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家：公会战斗次数
	plrcond_guildboss_fight_cnt: {

		gconst.Evt_GuildBossFight,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 玩家：竞技场排名(取负数)
	plrcond_arena_rank: {

		gconst.Evt_ArenaRank,

		func(p1 int32, cv ICondObj, args []interface{}) {
			rank := float64(args[1].(int32))
			if rank >= 0 {
				return
			}

			if cv.GetVal() == 0 || cv.GetVal() < rank {
				cv.SetVal(rank)
			}

			return
		},
	},

	// 玩家：开服庆典任务领取数
	plrcond_targetdays_cnt: {

		gconst.Evt_TargetDays_Take,

		func(p1 int32, cv ICondObj, args []interface{}) {
			tp := args[2].(int32)

			if p1 == -1 || p1 == tp {
				cv.AddVal(1)
			}

			return
		},
	},

	// 每月任务完成数量
	plrcond_taskmonth_cnt: {

		gconst.Evt_TaskMonth_Fin,

		func(p1 int32, cv ICondObj, args []interface{}) {
			tp := args[2].(int32)

			if p1 == -1 || p1 == tp {
				cv.AddVal(1)
			}

			return
		},
	},

	// 英灵试炼战斗次数
	plrcond_crusade_fight_cnt: {

		gconst.Evt_CrusadeFight,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 快速挂机次数
	plrcond_onekeygj_cnt: {

		gconst.Evt_WLevelOneKeyGj,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	// 爬塔扫荡次数
	plrcond_tower_raid_cnt: {

		gconst.Evt_TowerRaid,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_rift_monster_fight_cnt: {

		gconst.Evt_RiftMonsterFight,

		func(p1 int32, cv ICondObj, args []interface{}) {
			id := args[1].(int32)

			if p1 == -1 || p1 == id {
				cv.AddVal(1)
			}

			return
		},
	},

	plrcond_rift_mine_occupy_cnt: {

		gconst.Evt_RiftMineOccupy,

		func(p1 int32, cv ICondObj, args []interface{}) {
			id := args[1].(int32)

			if p1 == -1 || p1 == id {
				cv.AddVal(1)
			}

			return
		},
	},

	plrcond_rift_mine_take_cnt: {

		gconst.Evt_RiftMineTake,

		func(p1 int32, cv ICondObj, args []interface{}) {
			id := args[1].(int32)

			if p1 == -1 || p1 == id {
				cv.AddVal(1)
			}

			return
		},
	},

	plrcond_rift_box_take_cnt: {

		gconst.Evt_RiftBoxTake,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_ladder_fight_cnt: {

		gconst.Evt_LadderFight,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_guild_donate_cnt: {

		gconst.Evt_GuildDonate,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_guild_wish_cnt: {

		gconst.Evt_GuildWish,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_guild_help_cnt: {

		gconst.Evt_GuildHelp,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_guild_fin_order: {

		gconst.Evt_GuildOrderClose,

		func(p1 int32, cv ICondObj, args []interface{}) {
			star := args[1].(int32)

			if p1 == -1 || star == p1 {
				cv.AddVal(1)
			}

			return
		},
	},

	plrcond_guild_join_cnt: {

		gconst.Evt_GuildJoin,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_guild_sign_cnt: {

		gconst.Evt_GuildSign,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

		},
	},

	plrcond_appoint_send: {
		gconst.Evt_AppointSend,

		func(p1 int32, cv ICondObj, args []interface{}) {

			star := args[1].(int32)

			if p1 == -1 || star == p1 {
				cv.AddVal(1)
			}
			return
		},
	},

	plrcond_warcup_guess_cnt: {

		gconst.Evt_WarCupGuess,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

		},
	},

	plrcond_warcup_guess_win_cnt: {

		gconst.Evt_WarCupGuessWin,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_warcup_chat_cnt: {

		gconst.Evt_WarCupChat,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_warcup_watch_cnt: {

		gconst.Evt_WarCupWatch,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},

	plrcond_actmaze_cnt: {

		gconst.Evt_ActMaze_Click,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},
	plrcond_actmaze_event_cnt: {

		gconst.Evt_ActMaze_Event,

		func(p1 int32, cv ICondObj, args []interface{}) {

			cv.AddVal(1)

			return
		},
	},
	plrcond_actmaze_score: {

		gconst.Evt_ActMaze_Score,

		func(p1 int32, cv ICondObj, args []interface{}) {

			score := args[1].(int32)

			cv.AddVal(float64(score))

			return
		},
	},
}
