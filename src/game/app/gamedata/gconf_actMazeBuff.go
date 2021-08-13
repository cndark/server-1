package gamedata

var ConfActMazeBuff = &actMazeBuffTable{}

type actMazeBuff struct {
	BuffId     int32 `json:"buffId"`     // buffId
	ShowLvCond int32 `json:"showLvCond"` // 多少层才出现
	Weight     int32 `json:"weight"`     // 出现权重
}

type actMazeBuffTable struct {
	items map[int32]*actMazeBuff
}

func (self *actMazeBuffTable) Load() {
	var arr []*actMazeBuff
	if !load_json("actMazeBuff.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeBuff)

	for _, v := range arr {
		items[v.BuffId] = v
	}

	self.items = items
}

func (self *actMazeBuffTable) Query(buffId int32) *actMazeBuff {
	return self.items[buffId]
}

func (self *actMazeBuffTable) Items() map[int32]*actMazeBuff {
	return self.items
}
