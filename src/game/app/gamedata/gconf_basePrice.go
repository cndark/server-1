package gamedata

var ConfBasePrice = &basePriceTable{}

type basePrice struct {
	Id    int32 `json:"id"` // id（货品的id）
	Price []*struct {
		Id int32   `json:"id"`
		N  float64 `json:"n"`
	} `json:"price"` // 购买价格
}

type basePriceTable struct {
	items map[int32]*basePrice
}

func (self *basePriceTable) Load() {
	var arr []*basePrice
	if !load_json("basePrice.json", &arr) {
		return
	}

	items := make(map[int32]*basePrice)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *basePriceTable) Query(id int32) *basePrice {
	return self.items[id]
}

func (self *basePriceTable) Items() map[int32]*basePrice {
	return self.items
}
