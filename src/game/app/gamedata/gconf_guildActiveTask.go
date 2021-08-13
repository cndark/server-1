package gamedata

var ConfGuildActiveTask = &guildActiveTaskTable{}

type guildActiveTask struct {
	Id           int32   `json:"id"`           // 任务id
	ActiveReward int32   `json:"activeReward"` // 活跃度积分奖励
	Cond         int32   `json:"cond"`         // 任务条件ID
	P1           int32   `json:"p1"`           // 任务参数1
	P2           float64 `json:"p2"`           // 任务参数2
}

type guildActiveTaskTable struct {
	items map[int32]*guildActiveTask
}

func (self *guildActiveTaskTable) Load() {
	var arr []*guildActiveTask
	if !load_json("guildActiveTask.json", &arr) {
		return
	}

	items := make(map[int32]*guildActiveTask)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *guildActiveTaskTable) Query(id int32) *guildActiveTask {
	return self.items[id]
}

func (self *guildActiveTaskTable) Items() map[int32]*guildActiveTask {
	return self.items
}
