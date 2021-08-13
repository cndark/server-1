package gamedata

var ConfActMazeEventMine = &actMazeEventMineTable{}

type actMazeEventMine struct {
	Seq  int32 `json:"seq"` // 序列
	Drop []*struct {
		Item   int32 `json:"item"`
		DropId int32 `json:"dropId"`
	} `json:"drop"` // 消耗道具对应掉落库
	ShowLvCond int32 `json:"showLvCond"` // 多少层才出现
	Weight     int32 `json:"weight"`     // 出现权重
}

type actMazeEventMineTable struct {
	items map[int32]*actMazeEventMine
}

func (self *actMazeEventMineTable) Load() {
	var arr []*actMazeEventMine
	if !load_json("actMazeEventMine.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEventMine)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMazeEventMineTable) Query(seq int32) *actMazeEventMine {
	return self.items[seq]
}

func (self *actMazeEventMineTable) Items() map[int32]*actMazeEventMine {
	return self.items
}
