package gamedata

var ConfRiftMine = &riftMineTable{}

type riftMine struct {
	Id       int32   `json:"id"`       // id
	Txt_Name string  `json:"txt_Name"` // 名称
	Weight   int32   `json:"weight"`   // 权重
	TimeMin  []int32 `json:"timeMin"`  // 挖掘时间（分钟）
	Reward   []*struct {
		Id   int32   `json:"id"`
		N    int32   `json:"n"`
		Odds float32 `json:"odds"`
	} `json:"reward"` // 基础奖励值
}

type riftMineTable struct {
	items map[int32]*riftMine
}

func (self *riftMineTable) Load() {
	var arr []*riftMine
	if !load_json("riftMine.json", &arr) {
		return
	}

	items := make(map[int32]*riftMine)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *riftMineTable) Query(id int32) *riftMine {
	return self.items[id]
}

func (self *riftMineTable) Items() map[int32]*riftMine {
	return self.items
}
