package gamedata

var ConfActMonopolyTaskM = &actMonopolyTaskTableM{}

type actMonopolyTaskTableM struct {
	objs map[int32][]*actMonopolyTask
}

func (self *actMonopolyTaskTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyTask)

	for _, v := range ConfActMonopolyTask.Items() {
		key := v.ConfGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyTaskTableM) QueryItems(grp int32) []*actMonopolyTask {
	return self.objs[grp]
}
