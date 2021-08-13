package gamedata

var ConfActSumonOddsM = &actSummonOddsTableM{}

type actSummonOddsTableM struct {
	objs map[string][]*actSummonOdds
}

func (self *actSummonOddsTableM) Load() {
	self.objs = make(map[string][]*actSummonOdds)

	for _, v := range ConfActSummonOdds.Items() {
		key := v.SummonGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actSummonOddsTableM) Items(key string) []*actSummonOdds {
	return self.objs[key]
}
