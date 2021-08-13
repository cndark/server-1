package gamedata

var ConfActMazeEventMonsterM = &actMazeEventMonsterTableM{}

type actMazeEventMonsterTableM struct {
	objs map[int32][]*actMazeEventMonster
}

func (self *actMazeEventMonsterTableM) Load() {
	self.objs = make(map[int32][]*actMazeEventMonster)

	for _, v := range ConfActMazeEventMonster.Items() {
		key := v.Type
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMazeEventMonsterTableM) QueryItems(grp int32) []*actMazeEventMonster {
	return self.objs[grp]
}
