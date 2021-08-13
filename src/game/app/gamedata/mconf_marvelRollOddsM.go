package gamedata

var ConfMarvelRollOddsM = &marvelRollOddsM{}

type marvelRollOddsM struct {
	objs map[string]map[int32][]*marvelRollOdds
}

func (self *marvelRollOddsM) Load() {
	self.objs = make(map[string]map[int32][]*marvelRollOdds)

	grp_objs := map[string][]*marvelRollOdds{}
	for _, v := range ConfMarvelRollOdds.Items() {
		grp := v.Group

		grp_objs[grp] = append(grp_objs[grp], v)
	}

	for grp, v := range grp_objs {
		blank_objs := map[int32][]*marvelRollOdds{}

		for _, b := range v {
			bid := b.BlankId
			blank_objs[bid] = append(blank_objs[bid], b)
		}

		self.objs[grp] = blank_objs
	}
}

//返回，map key为格子id
func (self *marvelRollOddsM) QueryGrp(grp string) map[int32][]*marvelRollOdds {
	return self.objs[grp]
}

func (self *marvelRollOddsM) NewMarvelRollOddsToBlanks() map[int32][]*marvelRollOdds {
	return map[int32][]*marvelRollOdds{}
}
