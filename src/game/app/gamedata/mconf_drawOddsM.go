package gamedata

var ConfDrawOddsM = &drawOddsTableM{}

type drawOddsTableM struct {
	objs map[string][]*drawOdds
}

func (self *drawOddsTableM) Load() {
	self.objs = make(map[string][]*drawOdds)

	for _, v := range ConfDrawOdds.Items() {
		key := v.DrawGroup
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *drawOddsTableM) Items(key string) []*drawOdds {
	return self.objs[key]
}
