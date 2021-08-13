package gamedata

var ConfActRushLocalRewardM = &actRushLocalRewardTableM{}

type actRushLocalRewardTableM struct {
	objs map[int32][]*actRushLocalReward
}

func (self *actRushLocalRewardTableM) Load() {
	self.objs = make(map[int32][]*actRushLocalReward)

	for _, v := range ConfActRushLocalReward.Items() {
		key := v.Grp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actRushLocalRewardTableM) QueryItems(key int32) []*actRushLocalReward {
	return self.objs[key]
}
