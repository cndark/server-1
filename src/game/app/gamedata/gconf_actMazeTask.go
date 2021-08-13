package gamedata

var ConfActMazeTask = &actMazeTaskTable{}

type actMazeTask struct {
	Id     int32 `json:"id"`   // 任务id
	Type   int32 `json:"type"` // 任务类型，1每日，2成就
	Attain []*struct {
		AttainId int32   `json:"attainId"`
		P2       float64 `json:"p2"`
	} `json:"attain"` // 条件状态
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	ConfGrp int32 `json:"confGrp"` // 活动配置
}

type actMazeTaskTable struct {
	items map[int32]*actMazeTask
}

func (self *actMazeTaskTable) Load() {
	var arr []*actMazeTask
	if !load_json("actMazeTask.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeTask)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actMazeTaskTable) Query(id int32) *actMazeTask {
	return self.items[id]
}

func (self *actMazeTaskTable) Items() map[int32]*actMazeTask {
	return self.items
}
