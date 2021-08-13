package gamedata

var ConfWorldBossRankReward = &worldBossRankRewardTable{}

type worldBossRankReward struct {
	Id   int32 `json:"id"` // 序列
	Rank []*struct {
		Low  int32 `json:"low"`
		High int32 `json:"high"`
	} `json:"rank"` // 排名范围
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type worldBossRankRewardTable struct {
	items map[int32]*worldBossRankReward
}

func (self *worldBossRankRewardTable) Load() {
	var arr []*worldBossRankReward
	if !load_json("worldBossRankReward.json", &arr) {
		return
	}

	items := make(map[int32]*worldBossRankReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *worldBossRankRewardTable) Query(id int32) *worldBossRankReward {
	return self.items[id]
}

func (self *worldBossRankRewardTable) Items() map[int32]*worldBossRankReward {
	return self.items
}
