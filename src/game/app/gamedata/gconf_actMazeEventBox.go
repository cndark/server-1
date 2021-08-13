package gamedata

var ConfActMazeEventBox = &actMazeEventBoxTable{}

type actMazeEventBox struct {
	Seq  int32 `json:"seq"` // 序列
	Drop []*struct {
		Item   int32 `json:"item"`
		DropId int32 `json:"dropId"`
	} `json:"drop"` // 消耗道具对应掉落库
	ShowLvCond int32 `json:"showLvCond"` // 多少层才出现
	Weight     int32 `json:"weight"`     // 出现权重
}

type actMazeEventBoxTable struct {
	items map[int32]*actMazeEventBox
}

func (self *actMazeEventBoxTable) Load() {
	var arr []*actMazeEventBox
	if !load_json("actMazeEventBox.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEventBox)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMazeEventBoxTable) Query(seq int32) *actMazeEventBox {
	return self.items[seq]
}

func (self *actMazeEventBoxTable) Items() map[int32]*actMazeEventBox {
	return self.items
}
