package gamedata

var ConfActMonopolyProblem = &actMonopolyProblemTable{}

type actMonopolyProblem struct {
	Seq         int32 `json:"seq"`   // 序列
	Group       int32 `json:"group"` // 奖励组
	Right       int32 `json:"right"` // 正确答案
	RightReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"rightReward"` // 正确奖励
	ErrorReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"errorReward"` // 错误奖励
	Weight int32 `json:"weight"` // 出现权重
}

type actMonopolyProblemTable struct {
	items map[int32]*actMonopolyProblem
}

func (self *actMonopolyProblemTable) Load() {
	var arr []*actMonopolyProblem
	if !load_json("actMonopolyProblem.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyProblem)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMonopolyProblemTable) Query(seq int32) *actMonopolyProblem {
	return self.items[seq]
}

func (self *actMonopolyProblemTable) Items() map[int32]*actMonopolyProblem {
	return self.items
}
