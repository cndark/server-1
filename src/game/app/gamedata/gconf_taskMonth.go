package gamedata

var ConfTaskMonth = &taskMonthTable{}

type taskMonth struct {
	Id     int32   `json:"id"`   // 任务id
	Type   int32   `json:"type"` // 类型
	Cond   int32   `json:"cond"` // 条件
	P1     int32   `json:"p1"`   // 条件参数1
	P2     float64 `json:"p2"`   // 条件参数2
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type taskMonthTable struct {
	items map[int32]*taskMonth
}

func (self *taskMonthTable) Load() {
	var arr []*taskMonth
	if !load_json("taskMonth.json", &arr) {
		return
	}

	items := make(map[int32]*taskMonth)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *taskMonthTable) Query(id int32) *taskMonth {
	return self.items[id]
}

func (self *taskMonthTable) Items() map[int32]*taskMonth {
	return self.items
}
