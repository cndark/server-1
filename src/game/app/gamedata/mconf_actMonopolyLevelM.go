package gamedata

var ConfActMonopolyLevelM = &actMonopolyLevelTableM{}

type actMonopolyLevelTableM struct {
	objs map[int32][]*actMonopolyLevel
}

func (self *actMonopolyLevelTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyLevel)

	for _, v := range ConfActMonopolyLevel.Items() {
		key := v.ConfGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyLevelTableM) QueryItems(grp int32) []*actMonopolyLevel {
	return self.objs[grp]
}
