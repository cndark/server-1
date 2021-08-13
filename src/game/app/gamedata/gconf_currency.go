package gamedata

var ConfCurrency = &currencyTable{}

type currency struct {
	Id       int32  `json:"id"`       // 货币ID
	Txt_Name string `json:"txt_Name"` // 货币名称
}

type currencyTable struct {
	items map[int32]*currency
}

func (self *currencyTable) Load() {
	var arr []*currency
	if !load_json("currency.json", &arr) {
		return
	}

	items := make(map[int32]*currency)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *currencyTable) Query(id int32) *currency {
	return self.items[id]
}

func (self *currencyTable) Items() map[int32]*currency {
	return self.items
}
