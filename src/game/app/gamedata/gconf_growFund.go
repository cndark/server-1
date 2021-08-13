package gamedata

var ConfGrowFund = &growFundTable{}

type growFund struct {
	Id     int32 `json:"id"` // id
	Lv     int32 `json:"lv"` // 等级
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type growFundTable struct {
	items map[int32]*growFund
}

func (self *growFundTable) Load() {
	var arr []*growFund
	if !load_json("growFund.json", &arr) {
		return
	}

	items := make(map[int32]*growFund)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *growFundTable) Query(id int32) *growFund {
	return self.items[id]
}

func (self *growFundTable) Items() map[int32]*growFund {
	return self.items
}
