package utils

import "fw/src/game/app/gamedata"

// ============================================================================

// 公共翻牌奖励
func BattleReward(lv int32) map[int32]int64 {
	item := make(map[int32]int64)
	for _, conf := range gamedata.ConfBattleReward.Items() {
		if len(conf.LevelRange) > 0 &&
			lv >= conf.LevelRange[0].Low &&
			lv <= conf.LevelRange[0].High {

			slt := make(map[int32]int32)
			for idx, v := range conf.RewardGroup {
				slt[int32(idx+1)] += v.W
			}

			idx := PickWeightedMapId(slt)
			if idx != 0 {
				item[conf.RewardGroup[idx-1].Id] += int64(conf.RewardGroup[idx-1].N)
			}

			break
		}
	}

	return item
}
