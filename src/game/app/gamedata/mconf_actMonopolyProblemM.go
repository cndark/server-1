package gamedata

var ConfActMonopolyProblemM = &actMonopolyProblemTableM{}

type actMonopolyProblemTableM struct {
	objs map[int32][]*actMonopolyProblem
}

func (self *actMonopolyProblemTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyProblem)

	for _, v := range ConfActMonopolyProblem.Items() {
		key := v.Group
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyProblemTableM) QueryItems(grp int32) []*actMonopolyProblem {
	return self.objs[grp]
}
