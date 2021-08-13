package gamedata

var ConfCrusade = &crusadeTable{}

type crusade struct {
	Id           int32 `json:"id"` // 关卡id
	PowerCorrect []*struct {
		Low  float32 `json:"low"`
		High float32 `json:"high"`
	} `json:"powerCorrect"` // 战力匹配系数
	Reward []*struct {
		Low  int32 `json:"low"`
		High int32 `json:"high"`
		Id   int32 `json:"id"`
		N    int32 `json:"n"`
	} `json:"reward"` // 过关奖励
	Chest []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"chest"` // 宝箱
	RoundType int32 `json:"roundType"` // 回合类型
}

type crusadeTable struct {
	items map[int32]*crusade
}

func (self *crusadeTable) Load() {
	var arr []*crusade
	if !load_json("crusade.json", &arr) {
		return
	}

	items := make(map[int32]*crusade)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *crusadeTable) Query(id int32) *crusade {
	return self.items[id]
}

func (self *crusadeTable) Items() map[int32]*crusade {
	return self.items
}
