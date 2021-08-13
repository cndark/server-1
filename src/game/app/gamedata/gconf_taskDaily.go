package gamedata

var ConfTaskDaily = &taskDailyTable{}

type taskDaily struct {
	Id   int32   `json:"id"`   // 任务id
	Cond int32   `json:"cond"` // 条件
	P1   int32   `json:"p1"`   // 条件参数1
	P2   float64 `json:"p2"`   // 条件参数2
}

type taskDailyTable struct {
	items map[int32]*taskDaily
}

func (self *taskDailyTable) Load() {
	var arr []*taskDaily
	if !load_json("taskDaily.json", &arr) {
		return
	}

	items := make(map[int32]*taskDaily)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *taskDailyTable) Query(id int32) *taskDaily {
	return self.items[id]
}

func (self *taskDailyTable) Items() map[int32]*taskDaily {
	return self.items
}
