package gamedata

var ConfActMazeEventShop = &actMazeEventShopTable{}

type actMazeEventShop struct {
	Seq int32 `json:"seq"` // 序列
	Id  []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"id"` // 兑换物品id
	Cost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 兑换配方
	ExchangeLimit int32 `json:"exchangeLimit"` // 单次兑换上限
}

type actMazeEventShopTable struct {
	items map[int32]*actMazeEventShop
}

func (self *actMazeEventShopTable) Load() {
	var arr []*actMazeEventShop
	if !load_json("actMazeEventShop.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEventShop)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMazeEventShopTable) Query(seq int32) *actMazeEventShop {
	return self.items[seq]
}

func (self *actMazeEventShopTable) Items() map[int32]*actMazeEventShop {
	return self.items
}
