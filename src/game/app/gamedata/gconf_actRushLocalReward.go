package gamedata

var ConfActRushLocalReward = &actRushLocalRewardTable{}

type actRushLocalReward struct {
	Id   int32 `json:"id"`  // 序列
	Grp  int32 `json:"grp"` // 奖励组
	Rank []*struct {
		Low  int32 `json:"low"`
		High int32 `json:"high"`
	} `json:"rank"` // 排名范围
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type actRushLocalRewardTable struct {
	items map[int32]*actRushLocalReward
}

func (self *actRushLocalRewardTable) Load() {
	var arr []*actRushLocalReward
	if !load_json("actRushLocalReward.json", &arr) {
		return
	}

	items := make(map[int32]*actRushLocalReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actRushLocalRewardTable) Query(id int32) *actRushLocalReward {
	return self.items[id]
}

func (self *actRushLocalRewardTable) Items() map[int32]*actRushLocalReward {
	return self.items
}
