package gamedata

var ConfWarCupRankReward = &warCupRankRewardTable{}

type warCupRankReward struct {
	Id     int32 `json:"id"` // 序号
	Reward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
	WarCupRank string `json:"warCupRank"` // 邮件内容
}

type warCupRankRewardTable struct {
	items map[int32]*warCupRankReward
}

func (self *warCupRankRewardTable) Load() {
	var arr []*warCupRankReward
	if !load_json("warCupRankReward.json", &arr) {
		return
	}

	items := make(map[int32]*warCupRankReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *warCupRankRewardTable) Query(id int32) *warCupRankReward {
	return self.items[id]
}

func (self *warCupRankRewardTable) Items() map[int32]*warCupRankReward {
	return self.items
}
