package gamedata

var ConfBillRate = &billRateTable{}

type billRate struct {
	Key string  `json:"key"` // 源币种
	CNY float32 `json:"CNY"` // 人民币
	USD float32 `json:"USD"` // 美元
	KRW float32 `json:"KRW"` // 韩元
}

type billRateTable struct {
	items map[string]*billRate
}

func (self *billRateTable) Load() {
	var arr []*billRate
	if !load_json("billRate.json", &arr) {
		return
	}

	items := make(map[string]*billRate)

	for _, v := range arr {
		items[v.Key] = v
	}

	self.items = items
}

func (self *billRateTable) Query(key string) *billRate {
	return self.items[key]
}

func (self *billRateTable) Items() map[string]*billRate {
	return self.items
}
