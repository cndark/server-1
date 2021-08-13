package gamedata

var ConfBattleReward = &battleRewardTable{}

type battleReward struct {
	Id         int32 `json:"id"` // id
	LevelRange []*struct {
		Low  int32 `json:"low"`
		High int32 `json:"high"`
	} `json:"levelRange"` // 等级区间
	RewardGroup []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
		W  int32 `json:"w"`
	} `json:"rewardGroup"` // 奖励组
}

type battleRewardTable struct {
	items map[int32]*battleReward
}

func (self *battleRewardTable) Load() {
	var arr []*battleReward
	if !load_json("battleReward.json", &arr) {
		return
	}

	items := make(map[int32]*battleReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *battleRewardTable) Query(id int32) *battleReward {
	return self.items[id]
}

func (self *battleRewardTable) Items() map[int32]*battleReward {
	return self.items
}
