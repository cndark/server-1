package gamedata

var ConfShopItemM = &shopTableItemM{}

type shopTableItemM struct {
	objs map[int32]map[int32][]*shopItem
}

func (self *shopTableItemM) Load() {
	self.objs = make(map[int32]map[int32][]*shopItem)

	shop_objs := map[int32][]*shopItem{}
	for _, v := range ConfShopItem.Items() {
		shopid := v.ShopType

		shop_objs[shopid] = append(shop_objs[shopid], v)
	}

	for shopid, v := range shop_objs {
		shop_b_objs := map[int32][]*shopItem{}

		for _, s := range v {
			bid := s.BlankId
			shop_b_objs[bid] = append(shop_b_objs[bid], s)
		}

		self.objs[shopid] = shop_b_objs
	}

}

//返回单商店的物品，map key为格子id
func (self *shopTableItemM) QueryS(shopid int32) map[int32][]*shopItem {
	return self.objs[shopid]
}

func (self *shopTableItemM) NewShopItemToBlank() map[int32][]*shopItem {
	return map[int32][]*shopItem{}
}
