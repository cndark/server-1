package gamedata

var ConfActMazeTaskAttain = &actMazeTaskAttainTable{}

type actMazeTaskAttain struct {
	AttainId int32 `json:"attainId"` // 统计id
	Cond     int32 `json:"cond"`     // 条件id
	P1       int32 `json:"p1"`       // 条件参数
}

type actMazeTaskAttainTable struct {
	items map[int32]*actMazeTaskAttain
}

func (self *actMazeTaskAttainTable) Load() {
	var arr []*actMazeTaskAttain
	if !load_json("actMazeTaskAttain.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeTaskAttain)

	for _, v := range arr {
		items[v.AttainId] = v
	}

	self.items = items
}

func (self *actMazeTaskAttainTable) Query(attainId int32) *actMazeTaskAttain {
	return self.items[attainId]
}

func (self *actMazeTaskAttainTable) Items() map[int32]*actMazeTaskAttain {
	return self.items
}
