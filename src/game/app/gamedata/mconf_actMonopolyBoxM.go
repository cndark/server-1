package gamedata

var ConfActMonopolyBoxM = &actMonopolyBoxTableM{}

type actMonopolyBoxTableM struct {
	objs map[int32][]*actMonopolyBox
}

func (self *actMonopolyBoxTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyBox)

	for _, v := range ConfActMonopolyBox.Items() {
		key := v.Group
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyBoxTableM) QueryItems(grp int32) []*actMonopolyBox {
	return self.objs[grp]
}
