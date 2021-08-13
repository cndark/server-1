package gamedata

var ConfActRushLocalM = &actRushLocalTableM{}

type actRushLocalTableM struct {
	objs map[int32][]*actRushLocal
}

func (self *actRushLocalTableM) Load() {
	self.objs = make(map[int32][]*actRushLocal)

	for _, v := range ConfActRushLocal.Items() {
		key := v.ConfGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actRushLocalTableM) QueryItems(key int32) []*actRushLocal {
	return self.objs[key]
}
