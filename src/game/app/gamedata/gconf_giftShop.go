package gamedata

var ConfGiftShop = &giftShopTable{}

type giftShop struct {
	Id     int32 `json:"id"` // id
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	BuyCntLimit int32 `json:"buyCntLimit"` // 限购次数
	PayId       int32 `json:"payId"`       // payId
	Reset       int32 `json:"reset"`       // 重置，1每日/2每周/3每月
}

type giftShopTable struct {
	items map[int32]*giftShop
}

func (self *giftShopTable) Load() {
	var arr []*giftShop
	if !load_json("giftShop.json", &arr) {
		return
	}

	items := make(map[int32]*giftShop)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *giftShopTable) Query(id int32) *giftShop {
	return self.items[id]
}

func (self *giftShopTable) Items() map[int32]*giftShop {
	return self.items
}
