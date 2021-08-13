package gamedata

var ConfShop = &shopTable{}

type shop struct {
	Id          int32   `json:"Id"`          // 商城分类ID
	IsRandom    int32   `json:"isRandom"`    // 是否随机商品
	RefreshTime []int32 `json:"refreshTime"` // 刷新时间
	FreeCounter int32   `json:"freeCounter"` // 免费刷新
	RefreshCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"refreshCost"` // 付费刷新消耗
	OpenStatus []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"openStatus"` // 开启需要的状态
}

type shopTable struct {
	items map[int32]*shop
}

func (self *shopTable) Load() {
	var arr []*shop
	if !load_json("shop.json", &arr) {
		return
	}

	items := make(map[int32]*shop)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *shopTable) Query(Id int32) *shop {
	return self.items[Id]
}

func (self *shopTable) Items() map[int32]*shop {
	return self.items
}
