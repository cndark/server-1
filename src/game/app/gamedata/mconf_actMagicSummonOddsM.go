package gamedata

var ConfActMagicSummonOddsM = &actMagicSummonOddsTableM{}

type actMagicSummonOddsTableM struct {
	objs map[string][]*actMagicSummonOdds
}

func (self *actMagicSummonOddsTableM) Load() {
	self.objs = make(map[string][]*actMagicSummonOdds)

	for _, v := range ConfActMagicSummonOdds.Items() {
		key := v.SummonGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMagicSummonOddsTableM) Items(key string) []*actMagicSummonOdds {
	return self.objs[key]
}
