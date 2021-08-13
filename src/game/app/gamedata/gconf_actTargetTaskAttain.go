package gamedata

var ConfActTargetTaskAttain = &actTargetTaskAttainTable{}

type actTargetTaskAttain struct {
	AttainId int32 `json:"attainId"` // 统计id
	Cond     int32 `json:"cond"`     // 条件id
	P1       int32 `json:"p1"`       // 条件参数
}

type actTargetTaskAttainTable struct {
	items map[int32]*actTargetTaskAttain
}

func (self *actTargetTaskAttainTable) Load() {
	var arr []*actTargetTaskAttain
	if !load_json("actTargetTaskAttain.json", &arr) {
		return
	}

	items := make(map[int32]*actTargetTaskAttain)

	for _, v := range arr {
		items[v.AttainId] = v
	}

	self.items = items
}

func (self *actTargetTaskAttainTable) Query(attainId int32) *actTargetTaskAttain {
	return self.items[attainId]
}

func (self *actTargetTaskAttainTable) Items() map[int32]*actTargetTaskAttain {
	return self.items
}
