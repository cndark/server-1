package gamedata

var ConfDropGrpM = &dropGrpTableM{}

type dropGrpTableM struct {
	groups map[int32][]*dropGrp
	wsums  map[int32]int32
}

func (self *dropGrpTableM) Load() {
	self.groups = make(map[int32][]*dropGrp)
	self.wsums = make(map[int32]int32)

	for _, v := range ConfDropGrp.items {
		grp := v.GroupId

		self.groups[grp] = append(self.groups[grp], v)
		self.wsums[grp] += v.Weight
	}
}

func (self *dropGrpTableM) Query(grp int32) (arr []*dropGrp, wsum int32) {
	return self.groups[grp], self.wsums[grp]
}
