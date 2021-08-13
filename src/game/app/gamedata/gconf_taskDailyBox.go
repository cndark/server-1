package gamedata

var ConfTaskDailyBox = &taskDailyBoxTable{}

type taskDailyBox struct {
	Id         int32 `json:"id"`         // 宝箱序列
	ActiveCond int32 `json:"activeCond"` // 任务完成个数
	Reward     []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
}

type taskDailyBoxTable struct {
	items map[int32]*taskDailyBox
}

func (self *taskDailyBoxTable) Load() {
	var arr []*taskDailyBox
	if !load_json("taskDailyBox.json", &arr) {
		return
	}

	items := make(map[int32]*taskDailyBox)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *taskDailyBoxTable) Query(id int32) *taskDailyBox {
	return self.items[id]
}

func (self *taskDailyBoxTable) Items() map[int32]*taskDailyBox {
	return self.items
}
