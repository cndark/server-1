package gamedata

var ConfGuildRankReward = &guildRankRewardTable{}

type guildRankReward struct {
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

type guildRankRewardTable struct {
	items map[int32]*guildRankReward
}

func (self *guildRankRewardTable) Load() {
	var arr []*guildRankReward
	if !load_json("guildRankReward.json", &arr) {
		return
	}

	items := make(map[int32]*guildRankReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *guildRankRewardTable) Query(id int32) *guildRankReward {
	return self.items[id]
}

func (self *guildRankRewardTable) Items() map[int32]*guildRankReward {
	return self.items
}
