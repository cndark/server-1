package gamedata

var ConfActMazeEvent = &actMazeEventTable{}

type actMazeEvent struct {
	Id    int32 `json:"id"`    // 事件类型
	Score int32 `json:"score"` // 参与事件获得积分
}

type actMazeEventTable struct {
	items map[int32]*actMazeEvent
}

func (self *actMazeEventTable) Load() {
	var arr []*actMazeEvent
	if !load_json("actMazeEvent.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEvent)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actMazeEventTable) Query(id int32) *actMazeEvent {
	return self.items[id]
}

func (self *actMazeEventTable) Items() map[int32]*actMazeEvent {
	return self.items
}
