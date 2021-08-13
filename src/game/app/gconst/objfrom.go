package gconst

// ============================================================================
// 背包内容来源

const (
	ObjFrom_Init = 1 // init
	ObjFrom_GM   = 2 // GM
	ObjFrom_Sys  = 3 // sys

	ObjFrom_Bill_Normal = 100 // 钻石充值
	ObjFrom_Bill_Card   = 101 // 充值：月、周、终生卡
	ObjFrom_Bill_First  = 102 // 首冲项
	ObjFrom_Bill_Gift   = 103 // 充值礼包
	ObjFrom_Bill_Total  = 104 // 累计充值

	ObjFrom_GiftCode   = 200 // 礼包码
	ObjFrom_RefundCode = 201 // 返利码
	ObjFrom_GoldenHand = 202 // 点金手
	ObjFrom_OnlineBox  = 203 // 在线宝箱

	ObjFrom_Guild             = 300 // 家族
	ObjFrom_GuildCreate       = 301 // 家族创建
	ObjFrom_GuildChangeName   = 302 // 家族改名
	ObjFrom_GuildKick         = 303 // 弹劾族长
	ObjFrom_GuildSign         = 304 // 签到
	ObjFrom_GuildWishHelp     = 305 // 许愿助力
	ObjFrom_GuildWishRewards  = 306 // 许愿奖励
	ObjFrom_GuildOrderStarup  = 307 // 订单升星
	ObjFrom_GuildOrderClose   = 308 // 订单结束
	ObjFrom_GuildTechLevelup  = 309 // 科技升级
	ObjFrom_GuildTechReset    = 310 // 科技重置
	ObjFrom_GuildBossFight    = 311 // 打副本boss
	ObjFrom_GuildHarborDonate = 312 // 港口捐赠

	ObjFrom_CounterReset   = 400 // 计数器重置
	ObjFrom_CounterRecover = 401 // 计数器恢复

	ObjFrom_PlrChangeHFrame = 501 // 玩家改头像框
	ObjFrom_PlrChangeName   = 502 // 玩家改名
	ObjFrom_PlrLvUp         = 503 // 玩家升级
	ObjFrom_PlrOffline      = 505 // 离线奖励

	ObjFrom_Mail = 600 // 邮件

	ObjFrom_ItemUse      = 700 // 物品使用
	ObjFrom_ItemExchange = 701 // 物品快捷兑换
	ObjFrom_ItemChoose   = 702 // 物品选择

	ObjFrom_HeroLevelup        = 800 // 英雄升级
	ObjFrom_HeroStar           = 801 // 英雄升星
	ObjFrom_HeroReset          = 802 // 英雄重置
	ObjFrom_HeroDecompose      = 803 // 英雄分解
	ObjFrom_HeroChange         = 804 // 英雄转换
	ObjFrom_HeroTrinketUp      = 805 // 饰品升级
	ObjFrom_HeroTrinketRefresh = 806 // 饰品转换
	ObjFrom_HeroBagBuy         = 807 // 英雄背包购买
	ObjFrom_HeroInherit        = 808 // 英雄继承
	ObjFrom_HeroSkin           = 809 // 英雄皮肤

	ObjFrom_ArmorCompose = 900 // 装备合成

	ObjFrom_RelicEat = 1000 // 神器吃

	ObjFrom_RankPlayLike = 1100 // 榜单点赞

	ObjFrom_GWarFight = 1200 // 公会战打架

	ObjFrom_TaskMain  = 1600 // 主线任务
	ObjFrom_TaskDaily = 1601 // 日常任务
	ObjFrom_TaskAchv  = 1602 // 成就任务
	ObjFrom_TaskMonth = 1603 // 每月任务

	ObjFrom_TargetDays = 1700 // 开服庆典任务

	ObjFrom_Chat = 2700 // 聊天

	ObjFrom_ActRushLocal   = 3100 // 活动: 限时冲榜
	ObjFrom_ActBillLtTotal = 3102 // 活动: 限时累计充值
	ObjFrom_ActBillLtDay   = 3103 // 活动：限时累天充值
	ObjFrom_ActGift        = 3104 // 活动：礼包
	ObjFrom_ActSummon      = 3105 // 活动：主题活动
	ObjFrom_ActTargetTask  = 3106 // 活动：达标任务
	ObjFrom_ActMagicSummon = 3107 // 活动：魔法召唤
	ObjFrom_ActMonopoly    = 3109 // 活动：大富翁
	ObjFrom_ActMaze        = 3110 // 活动：迷宫

	ObjFrom_MOpen     = 3700 // open
	ObjFrom_MOpenTask = 3701 // 预告
	ObjFrom_MOpenBook = 3702 // 预约

	ObjFrom_AIBot = 4200 // 机器人

	ObjFrom_Draw    = 4300 // 抽卡
	ObjFrom_DrawBox = 4301 // 抽卡宝箱

	ObjFrom_WLevelFight       = 4400 // 推图战斗
	ObjFrom_WLevelGJ          = 4401 // 挂机奖励
	ObjFrom_WLevelOneKeyGJ    = 4402 // 一键挂机
	ObjFrom_WLevelDraw        = 4403 // 推图十连抽卡
	ObjFrom_WLevelOneKeyFight = 4404 // 一键推图

	ObjFrom_AppointRefresh = 4500 // 酒馆委派刷新
	ObjFrom_AppointSend    = 4501 // 酒馆委派
	ObjFrom_AppointAcc     = 4502 // 酒馆委派加速
	ObjFrom_AppointTake    = 4503 // 酒馆委派领奖

	ObjFrom_TowerFight = 4600 // 爬塔战斗
	ObjFrom_TowerRaid  = 4601 // 爬塔扫荡

	ObjFrom_ArenaFight = 4700 // 竞技场比武

	ObjFrom_MarvelRollRefresh = 4800 // 奇迹转盘-刷新
	ObjFrom_MarvelRollTake    = 4801 // 奇迹转盘-转

	ObjFrom_FriendGive = 4901 // 好友送礼

	ObjFrom_CrusadeBox   = 5001 // 远征宝箱
	ObjFrom_CrusadeFight = 5002 // 远征打架

	ObjFrom_Vip = 5004 // "vip"

	ObjFrom_MonthTicket = 5100 // 月票
	ObjFrom_PushGift    = 5101 // 推送礼包
	ObjFrom_GiftShop    = 5102 // 礼包商店
	ObjFrom_PrivCard    = 5103 // 特权卡
	ObjFrom_SignDaily   = 5104 // 每日签到
	ObjFrom_WLevelFund  = 5105 // 关卡基金
	ObjFrom_GrowFund    = 5106 // 成长基金

	ObjFrom_DaySign   = 5210 // 七日之约
	ObjFrom_TaskGrow  = 5211 // 进阶之路
	ObjFrom_BillFirst = 5212 // 超值首充

	ObjFrom_RiftMonster = 5300 // 裂隙怪物
	ObjFrom_RiftMine    = 5301 // 裂隙矿
	ObjFrom_RiftBox     = 5302 // 裂隙宝箱

	ObjFrom_LadderFight = 5400 // 天梯打架

	ObjFrom_WarCup = 5501 // 本服杯赛

	ObjFrom_WBossFight     = 5600 // 世界老板打架
	ObjFrom_WBossMaxDmgRwd = 5601 // 世界老板最大伤害奖

	ObjFrom_Invite = 5701 // 分享收藏

	ObjFrom_ShopBuy     = 10000 // 商店购买
	ObjFrom_ShopRefresh = 20000 // 商店刷新

	// 1000001 往后给fields counter占用
)
