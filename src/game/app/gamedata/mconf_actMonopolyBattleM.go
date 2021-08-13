package gamedata

var ConfActMonopolyBattleM = &actMonopolyBattleTableM{}

type actMonopolyBattleTableM struct {
	objs map[int32][]*actMonopolyBattle
}

func (self *actMonopolyBattleTableM) Load() {
	self.objs = make(map[int32][]*actMonopolyBattle)

	for _, v := range ConfActMonopolyBattle.Items() {
		key := v.Group
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMonopolyBattleTableM) QueryItems(grp int32) []*actMonopolyBattle {
	return self.objs[grp]
}
