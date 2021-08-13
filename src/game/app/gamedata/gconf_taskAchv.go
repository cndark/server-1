package gamedata

var ConfTaskAchv = &taskAchvTable{}

type taskAchv struct {
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

type taskAchvTable struct {
	items map[int32]*taskAchv
}

func (self *taskAchvTable) Load() {
	var arr []*taskAchv
	if !load_json("taskAchv.json", &arr) {
		return
	}

	items := make(map[int32]*taskAchv)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *taskAchvTable) Query(id int32) *taskAchv {
	return self.items[id]
}

func (self *taskAchvTable) Items() map[int32]*taskAchv {
	return self.items
}
