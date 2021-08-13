package gamedata

var ConfMonthTicketTask = &monthTicketTaskTable{}

type monthTicketTask struct {
	Id     int32   `json:"id"`     // id
	GetExp int32   `json:"getExp"` // 获得经验
	Cond   int32   `json:"cond"`   // 条件
	P1     int32   `json:"p1"`     // p1
	P2     float64 `json:"p2"`     // p2
}

type monthTicketTaskTable struct {
	items map[int32]*monthTicketTask
}

func (self *monthTicketTaskTable) Load() {
	var arr []*monthTicketTask
	if !load_json("monthTicketTask.json", &arr) {
		return
	}

	items := make(map[int32]*monthTicketTask)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *monthTicketTaskTable) Query(id int32) *monthTicketTask {
	return self.items[id]
}

func (self *monthTicketTaskTable) Items() map[int32]*monthTicketTask {
	return self.items
}
