package gamedata

var ConfGrowFundReward = &growFundRewardTable{}

type growFundReward struct {
	Id     int32 `json:"id"`     // id
	Lv     int32 `json:"lv"`     // 等级
	BuyCnt int32 `json:"buyCnt"` // 购买人数
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type growFundRewardTable struct {
	items map[int32]*growFundReward
}

func (self *growFundRewardTable) Load() {
	var arr []*growFundReward
	if !load_json("growFundReward.json", &arr) {
		return
	}

	items := make(map[int32]*growFundReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *growFundRewardTable) Query(id int32) *growFundReward {
	return self.items[id]
}

func (self *growFundRewardTable) Items() map[int32]*growFundReward {
	return self.items
}
