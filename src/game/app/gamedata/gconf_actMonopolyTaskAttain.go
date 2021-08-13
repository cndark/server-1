package gamedata

var ConfActMonopolyTaskAttain = &actMonopolyTaskAttainTable{}

type actMonopolyTaskAttain struct {
	AttainId int32 `json:"attainId"` // 统计id
	Cond     int32 `json:"cond"`     // 条件id
	P1       int32 `json:"p1"`       // 条件参数
}

type actMonopolyTaskAttainTable struct {
	items map[int32]*actMonopolyTaskAttain
}

func (self *actMonopolyTaskAttainTable) Load() {
	var arr []*actMonopolyTaskAttain
	if !load_json("actMonopolyTaskAttain.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyTaskAttain)

	for _, v := range arr {
		items[v.AttainId] = v
	}

	self.items = items
}

func (self *actMonopolyTaskAttainTable) Query(attainId int32) *actMonopolyTaskAttain {
	return self.items[attainId]
}

func (self *actMonopolyTaskAttainTable) Items() map[int32]*actMonopolyTaskAttain {
	return self.items
}
