package gamedata

var ConfRiftMonster = &riftMonsterTable{}

type riftMonster struct {
	Id       int32   `json:"id"`       // id
	Txt_Name string  `json:"txt_Name"` // 名称
	Weight   int32   `json:"weight"`   // 权重
	Round    int32   `json:"round"`    // 战斗回合
	Monster  []int32 `json:"monster"`  // 关卡怪物
	Reward   []*struct {
		Id   int32   `json:"id"`
		N    int32   `json:"n"`
		Odds float32 `json:"odds"`
	} `json:"reward"` // 通关奖励
}

type riftMonsterTable struct {
	items map[int32]*riftMonster
}

func (self *riftMonsterTable) Load() {
	var arr []*riftMonster
	if !load_json("riftMonster.json", &arr) {
		return
	}

	items := make(map[int32]*riftMonster)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *riftMonsterTable) Query(id int32) *riftMonster {
	return self.items[id]
}

func (self *riftMonsterTable) Items() map[int32]*riftMonster {
	return self.items
}
