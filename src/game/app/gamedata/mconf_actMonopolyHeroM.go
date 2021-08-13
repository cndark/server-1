package gamedata

var ConfActMonopolyHeroM = &actMonopolyHeroTableM{}

type actMonopolyHeroTableM struct {
	objs map[int32][]*actMonopolyHero
}

func (self *actMonopolyHeroTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyHero)

	for _, v := range ConfActMonopolyHero.Items() {
		key := v.Group
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyHeroTableM) QueryItems(grp int32) []*actMonopolyHero {
	return self.objs[grp]
}
