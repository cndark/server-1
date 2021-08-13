package gamedata

var ConfStaffRankReward = &staffRankRewardTable{}

type staffRankReward struct {
	Rank   int32 `json:"rank"` // 排名
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type staffRankRewardTable struct {
	items map[int32]*staffRankReward
}

func (self *staffRankRewardTable) Load() {
	var arr []*staffRankReward
	if !load_json("staffRankReward.json", &arr) {
		return
	}

	items := make(map[int32]*staffRankReward)

	for _, v := range arr {
		items[v.Rank] = v
	}

	self.items = items
}

func (self *staffRankRewardTable) Query(rank int32) *staffRankReward {
	return self.items[rank]
}

func (self *staffRankRewardTable) Items() map[int32]*staffRankReward {
	return self.items
}
