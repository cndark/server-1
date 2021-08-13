package gamedata

var ConfActMazeRankM = &actMazeRankTableM{}

type actMazeRankTableM struct {
	objs map[int32][]*actMazeRank
}

func (self *actMazeRankTableM) Load() {
	self.objs = make(map[int32][]*actMazeRank)

	for _, v := range ConfActMazeRank.Items() {
		key := v.ConfGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actMazeRankTableM) QueryItems(grp int32) []*actMazeRank {
	return self.objs[grp]
}
