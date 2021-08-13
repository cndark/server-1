package gamedata

var ConfActMonopolyShopM = &actMonopolyShopTableM{}

type actMonopolyShopTableM struct {
	objs map[int32][]*actMonopolyShop
}

func (self *actMonopolyShopTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyShop)

	for _, v := range ConfActMonopolyShop.Items() {
		key := v.Group
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyShopTableM) QueryItems(grp int32) []*actMonopolyShop {
	return self.objs[grp]
}
