package gamedata

var ConfActMonopolyTask = &actMonopolyTaskTable{}

type actMonopolyTask struct {
	Id     int32 `json:"id"` // 任务id
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

type actMonopolyTaskTable struct {
	items map[int32]*actMonopolyTask
}

func (self *actMonopolyTaskTable) Load() {
	var arr []*actMonopolyTask
	if !load_json("actMonopolyTask.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyTask)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actMonopolyTaskTable) Query(id int32) *actMonopolyTask {
	return self.items[id]
}

func (self *actMonopolyTaskTable) Items() map[int32]*actMonopolyTask {
	return self.items
}
