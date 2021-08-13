package gamedata

var ConfActTargetTask = &actTargetTaskTable{}

type actTargetTask struct {
	Id     int32 `json:"id"` // 任务id
	Attain []*struct {
		AttainId int32   `json:"attainId"`
		P2       float64 `json:"p2"`
	} `json:"attain"` // 条件状态
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	IsReset int32 `json:"isReset"` // 是否可重置
	ConfGrp int32 `json:"confGrp"` // 活动配置
}

type actTargetTaskTable struct {
	items map[int32]*actTargetTask
}

func (self *actTargetTaskTable) Load() {
	var arr []*actTargetTask
	if !load_json("actTargetTask.json", &arr) {
		return
	}

	items := make(map[int32]*actTargetTask)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actTargetTaskTable) Query(id int32) *actTargetTask {
	return self.items[id]
}

func (self *actTargetTaskTable) Items() map[int32]*actTargetTask {
	return self.items
}
