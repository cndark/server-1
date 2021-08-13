package gamedata

var ConfItemExchange = &itemExchangeTable{}

type itemExchange struct {
	Id   int32 `json:"id"` // 物品id
	Cost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"cost"` // 购买消耗
}

type itemExchangeTable struct {
	items map[int32]*itemExchange
}

func (self *itemExchangeTable) Load() {
	var arr []*itemExchange
	if !load_json("itemExchange.json", &arr) {
		return
	}

	items := make(map[int32]*itemExchange)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *itemExchangeTable) Query(id int32) *itemExchange {
	return self.items[id]
}

func (self *itemExchangeTable) Items() map[int32]*itemExchange {
	return self.items
}
