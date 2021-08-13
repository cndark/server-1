package gconst

// ============================================================================
// 排行榜 Id

const (
	RankId_AtkPwr   = 1 // 玩家战力榜
	RankId_WLevelLv = 2 // 玩家推图榜
	RankId_TowerLv  = 3 // 玩家爬塔榜

	// --------------------------------
	// 限时冲榜--本服 [101-199]

	RankId_RushLocal_AtkPwr       = 101 // 战力冲榜
	RankId_RushLocal_WLevel       = 102 // 主线推图榜
	RankId_RushLocal_Draw         = 103 // 召唤冲榜
	RankId_RushLocal_Arena        = 104 // 比武场战斗冲榜
	RankId_RushLocal_Tower        = 105 // 试炼冲榜
	RankId_RushLocal_GuildBossDmg = 106 // 家族boss伤害榜
	RankId_RushLocal_MarvelRoll   = 107 // 奇迹之盘冲榜

	// -------------------------------
	// 活动榜单 [301-399]
	RankId_ActMaze_Score = 301 // 迷宫积分排行榜

)

func IsGuildRankType(rk int32) bool {
	return false
}

func RankMaxRows(rkid int32) int {

	// 本服冲榜
	if rkid == RankId_RushLocal_WLevel {
		return 1000
	} else if rkid >= RankId_RushLocal_AtkPwr &&
		rkid <= RankId_RushLocal_MarvelRoll {
		return 500
	}

	return 200
}
