package gamedata

var ConfActMazeEventGrp = &actMazeEventGrpTable{}

type actMazeEventGrp struct {
	Grp      int32 `json:"grp"` // 组合id
	EventNum []*struct {
		EventId int32 `json:"eventId"`
		N       int32 `json:"n"`
	} `json:"eventNum"` // 事件数量
}

type actMazeEventGrpTable struct {
	items map[int32]*actMazeEventGrp
}

func (self *actMazeEventGrpTable) Load() {
	var arr []*actMazeEventGrp
	if !load_json("actMazeEventGrp.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEventGrp)

	for _, v := range arr {
		items[v.Grp] = v
	}

	self.items = items
}

func (self *actMazeEventGrpTable) Query(grp int32) *actMazeEventGrp {
	return self.items[grp]
}

func (self *actMazeEventGrpTable) Items() map[int32]*actMazeEventGrp {
	return self.items
}
