package gamedata

var ConfActMazeEventItem = &actMazeEventItemTable{}

type actMazeEventItem struct {
	Seq    int32 `json:"seq"` // 序列
	Reward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 道具奖励
	ShowLvCond int32 `json:"showLvCond"` // 多少层才出现
	Weight     int32 `json:"weight"`     // 出现权重
}

type actMazeEventItemTable struct {
	items map[int32]*actMazeEventItem
}

func (self *actMazeEventItemTable) Load() {
	var arr []*actMazeEventItem
	if !load_json("actMazeEventItem.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEventItem)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMazeEventItemTable) Query(seq int32) *actMazeEventItem {
	return self.items[seq]
}

func (self *actMazeEventItemTable) Items() map[int32]*actMazeEventItem {
	return self.items
}
