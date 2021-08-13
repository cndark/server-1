package gamedata

var ConfCounter = &counterTable{}

type counter struct {
	Id       int32 `json:"id"`       // 计数器id
	MaxValue int32 `json:"maxValue"` // 初始上限值
	Type     int32 `json:"type"`     // 恢复类型
	Recover  []*struct {
		Sec int32 `json:"sec"`
		N   int32 `json:"n"`
	} `json:"recover"` // 恢复值
	BuyId  int32 `json:"buyId"` // 购买消耗的id(别的counter)
	BuyN   int64 `json:"buyN"`  // 每次购买获得几点(0表示补满)
	OpCost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"opCost"` // 计数变化的价格(自身消耗顺带消耗的道具）
}

type counterTable struct {
	items map[int32]*counter
}

func (self *counterTable) Load() {
	var arr []*counter
	if !load_json("counter.json", &arr) {
		return
	}

	items := make(map[int32]*counter)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *counterTable) Query(id int32) *counter {
	return self.items[id]
}

func (self *counterTable) Items() map[int32]*counter {
	return self.items
}
