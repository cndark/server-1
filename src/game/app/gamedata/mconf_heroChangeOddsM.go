package gamedata

var ConfHeroChangeOddsM = &heroChangeOddsTableM{}

type heroChangeOddsTableM struct {
	objs map[int32][]*heroChangeOdds
}

func (self *heroChangeOddsTableM) Load() {
	self.objs = make(map[int32][]*heroChangeOdds)

	for _, v := range ConfHeroChangeOdds.Items() {
		grp := v.Group
		self.objs[grp] = append(self.objs[grp], v)
	}
}

func (self *heroChangeOddsTableM) Items(grp int32) []*heroChangeOdds {
	return self.objs[grp]
}
