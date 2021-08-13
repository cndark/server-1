package gamedata

var ConfGwarRankReward = &gwarRankRewardTable{}

type gwarRankReward struct {
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

type gwarRankRewardTable struct {
	items map[int32]*gwarRankReward
}

func (self *gwarRankRewardTable) Load() {
	var arr []*gwarRankReward
	if !load_json("gwarRankReward.json", &arr) {
		return
	}

	items := make(map[int32]*gwarRankReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *gwarRankRewardTable) Query(id int32) *gwarRankReward {
	return self.items[id]
}

func (self *gwarRankRewardTable) Items() map[int32]*gwarRankReward {
	return self.items
}
