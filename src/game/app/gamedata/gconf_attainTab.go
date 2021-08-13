package gamedata

var ConfAttainTab = &attainTabTable{}

type attainTab struct {
	AttainId int32 `json:"attainId"` // 统计id
	Cond     int32 `json:"cond"`     // 条件id
	P1       int32 `json:"p1"`       // 条件参数
}

type attainTabTable struct {
	items map[int32]*attainTab
}

func (self *attainTabTable) Load() {
	var arr []*attainTab
	if !load_json("attainTab.json", &arr) {
		return
	}

	items := make(map[int32]*attainTab)

	for _, v := range arr {
		items[v.AttainId] = v
	}

	self.items = items
}

func (self *attainTabTable) Query(attainId int32) *attainTab {
	return self.items[attainId]
}

func (self *attainTabTable) Items() map[int32]*attainTab {
	return self.items
}
