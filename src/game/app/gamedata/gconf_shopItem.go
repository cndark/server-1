package gamedata

var ConfShopItem = &shopItemTable{}

type shopItem struct {
	Id        int32   `json:"Id"`       // 商店物品生成序列\n（id=商店类型*100000+格子id*1000+自然序列）
	ShopType  int32   `json:"shopType"` // 商店类型
	BlankId   int32   `json:"blankId"`  // 格子ID
	Item      int32   `json:"item"`     // 物品配置
	Num       int32   `json:"num"`      // 数量
	Odds      int32   `json:"odds"`     // 出现权重
	Discount  int32   `json:"discount"` // 物品折价率
	Currency  int32   `json:"currency"` // 需要的货币
	AddPrice  float64 `json:"addPrice"` // 购买价格递增
	BuyCount  int32   `json:"buyCount"` // 可购买次数
	BlankOpen []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"blankOpen"` // 格子开启条件
	ItemOpen []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"itemOpen"` // 商品开启条件
	ItemClose []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"itemClose"` // 商品关闭条件
}

type shopItemTable struct {
	items map[int32]*shopItem
}

func (self *shopItemTable) Load() {
	var arr []*shopItem
	if !load_json("shopItem.json", &arr) {
		return
	}

	items := make(map[int32]*shopItem)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *shopItemTable) Query(Id int32) *shopItem {
	return self.items[Id]
}

func (self *shopItemTable) Items() map[int32]*shopItem {
	return self.items
}
