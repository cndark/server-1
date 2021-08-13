package gamedata

var ConfActMonopolyShop = &actMonopolyShopTable{}

type actMonopolyShop struct {
	Seq      int32 `json:"seq"`   // 序列
	Group    int32 `json:"group"` // 奖励组
	ShopItem []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"shopItem"` // 商品
	Ccy             int32 `json:"ccy"` // 货币Id
	DiscountAndOdds []*struct {
		Weight   int32 `json:"weight"`
		Discount int32 `json:"discount"`
	} `json:"discountAndOdds"` // 折扣和权重
	Weight int32 `json:"weight"` // 商品出现权重
}

type actMonopolyShopTable struct {
	items map[int32]*actMonopolyShop
}

func (self *actMonopolyShopTable) Load() {
	var arr []*actMonopolyShop
	if !load_json("actMonopolyShop.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyShop)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMonopolyShopTable) Query(seq int32) *actMonopolyShop {
	return self.items[seq]
}

func (self *actMonopolyShopTable) Items() map[int32]*actMonopolyShop {
	return self.items
}
