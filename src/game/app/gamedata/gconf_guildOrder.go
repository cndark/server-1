package gamedata

var ConfGuildOrder = &guildOrderTable{}

type guildOrder struct {
	Star   int32 `json:"star"`   // 星级
	Weight int32 `json:"weight"` // 任务权重
	Time   int32 `json:"time"`   // 探索时间(分钟)
	Reward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
	UpStarCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"upStarCost"` // 升级消耗
}

type guildOrderTable struct {
	items map[int32]*guildOrder
}

func (self *guildOrderTable) Load() {
	var arr []*guildOrder
	if !load_json("guildOrder.json", &arr) {
		return
	}

	items := make(map[int32]*guildOrder)

	for _, v := range arr {
		items[v.Star] = v
	}

	self.items = items
}

func (self *guildOrderTable) Query(star int32) *guildOrder {
	return self.items[star]
}

func (self *guildOrderTable) Items() map[int32]*guildOrder {
	return self.items
}
