package gamedata

var ConfWorldLevelFund = &worldLevelFundTable{}

type worldLevelFund struct {
	Id     int32 `json:"id"` // id
	Lv     int32 `json:"lv"` // 关卡
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type worldLevelFundTable struct {
	items map[int32]*worldLevelFund
}

func (self *worldLevelFundTable) Load() {
	var arr []*worldLevelFund
	if !load_json("worldLevelFund.json", &arr) {
		return
	}

	items := make(map[int32]*worldLevelFund)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *worldLevelFundTable) Query(id int32) *worldLevelFund {
	return self.items[id]
}

func (self *worldLevelFundTable) Items() map[int32]*worldLevelFund {
	return self.items
}
