package gamedata

var ConfWarCupGuessTask = &warCupGuessTaskTable{}

type warCupGuessTask struct {
	Id              int32 `json:"id"` // 任务id
	GuessTaskAttain []*struct {
		AttainId int32   `json:"attainId"`
		P2       float64 `json:"p2"`
	} `json:"guessTaskAttain"` // 统计条件
	Reward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
	ConfGrp int32 `json:"confGrp"` // 奖励组
}

type warCupGuessTaskTable struct {
	items map[int32]*warCupGuessTask
}

func (self *warCupGuessTaskTable) Load() {
	var arr []*warCupGuessTask
	if !load_json("warCupGuessTask.json", &arr) {
		return
	}

	items := make(map[int32]*warCupGuessTask)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *warCupGuessTaskTable) Query(id int32) *warCupGuessTask {
	return self.items[id]
}

func (self *warCupGuessTaskTable) Items() map[int32]*warCupGuessTask {
	return self.items
}
