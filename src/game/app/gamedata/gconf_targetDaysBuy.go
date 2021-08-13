package gamedata

var ConfTargetDaysBuy = &targetDaysBuyTable{}

type targetDaysBuy struct {
	Id   int32 `json:"id"`  // id
	Day  int32 `json:"day"` // 天数
	Item []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"item"` // 物品
	Price  int32 `json:"price"`  // 现价
	MaxBuy int32 `json:"maxBuy"` // 限购
}

type targetDaysBuyTable struct {
	items map[int32]*targetDaysBuy
}

func (self *targetDaysBuyTable) Load() {
	var arr []*targetDaysBuy
	if !load_json("targetDaysBuy.json", &arr) {
		return
	}

	items := make(map[int32]*targetDaysBuy)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *targetDaysBuyTable) Query(id int32) *targetDaysBuy {
	return self.items[id]
}

func (self *targetDaysBuyTable) Items() map[int32]*targetDaysBuy {
	return self.items
}
