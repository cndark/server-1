package gconst

// ============================================================================

// 玩家计数器
// 100000x 正常计数(相互独立)
// 110000x 对应购买100000x的计数(需要有100000x成对出现)

const (
	Cnt_WLevelOneKeyGj    = 1000001 // 一键挂机次数
	Cnt_AppointTask       = 1000002 // 酒馆派遣数量(投放,仅使用上限值)
	Cnt_TowerFightBuy     = 1000003 // 爬塔挑战体力购买次数
	Cnt_TowerFight        = 1100003 // 爬塔挑战体力数
	Cnt_HeroBagBuy        = 1000004 // 英雄背包扩容次数
	Cnt_HeroBag           = 1000005 // 英雄背包初始容量数(投放,仅使用上限值)
	Cnt_PlayerStrengthBuy = 1000006 // 购买玩家体力次数
	Cnt_PlayerStrength    = 1100006 // 玩家体力数
	Cnt_GoldenHand        = 1000007 // 点金手
	Cnt_WLevelGJTime      = 1000010 // 挂机上限额外时间(投放,仅使用上限值)
	Cnt_WLevelGJExtReward = 1000011 // 挂机额外奖励(投放,仅使用上限值)
	Cnt_WLevelFightOneKey = 1000012 // 快速战斗次数
	Cnt_DrawNormal_Free   = 1000013 // 普通召唤免费次数上限

	Cnt_GuildWishHelp  = 1000101 // 公会祈愿助力次数
	Cnt_GuildTechReset = 1000102 // 公会科技重置
	Cnt_GuildBossFight = 1000103 // 公会boss挑战
	Cnt_GWarFight      = 1000104 // 公会战挑战
	Cnt_RankLikeWLevel = 1000201 // 榜单点赞推图
	Cnt_RankLikeTower  = 1000202 // 榜单点赞爬塔
	Cnt_WBossFight     = 1100301 // 世界boss打架次数
	Cnt_MazeBuyCnt     = 1000311 // 迷宫体力购买次数
	Cnt_MazeCnt        = 1100311 // 迷宫体力
	Cnt_MazeBuyLife    = 1000312 // 迷宫血量购买次数
	Cnt_MazeLife       = 1100312 // 迷宫血量
)

// 计数器回复类型

const (
	CounterType_Unlimit   = 0 // 不恢复，无上限限制
	CounterType_UnRecover = 1 // 不恢复，有上限限制
	CounterType_Recover   = 2 // 时间恢复，有上限限制
	CounterType_Daily     = 3 // 每天重置，有上限限制
	CounterType_Weekly    = 4 // 每周重置，有上限限制
)
