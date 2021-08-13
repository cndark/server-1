package gamedata

var ConfTaskGrow = &taskGrowTable{}

type taskGrow struct {
	Id        int32 `json:"id"` // 任务id
	AttainTab []*struct {
		AttainId int32   `json:"attainId"`
		P2       float64 `json:"p2"`
	} `json:"attainTab"` // 条件状态
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type taskGrowTable struct {
	items map[int32]*taskGrow
}

func (self *taskGrowTable) Load() {
	var arr []*taskGrow
	if !load_json("taskGrow.json", &arr) {
		return
	}

	items := make(map[int32]*taskGrow)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *taskGrowTable) Query(id int32) *taskGrow {
	return self.items[id]
}

func (self *taskGrowTable) Items() map[int32]*taskGrow {
	return self.items
}
