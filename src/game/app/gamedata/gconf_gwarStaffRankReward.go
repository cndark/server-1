package gamedata

var ConfGwarStaffRankReward = &gwarStaffRankRewardTable{}

type gwarStaffRankReward struct {
	Rank   int32 `json:"rank"` // 排名
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type gwarStaffRankRewardTable struct {
	items map[int32]*gwarStaffRankReward
}

func (self *gwarStaffRankRewardTable) Load() {
	var arr []*gwarStaffRankReward
	if !load_json("gwarStaffRankReward.json", &arr) {
		return
	}

	items := make(map[int32]*gwarStaffRankReward)

	for _, v := range arr {
		items[v.Rank] = v
	}

	self.items = items
}

func (self *gwarStaffRankRewardTable) Query(rank int32) *gwarStaffRankReward {
	return self.items[rank]
}

func (self *gwarStaffRankRewardTable) Items() map[int32]*gwarStaffRankReward {
	return self.items
}
