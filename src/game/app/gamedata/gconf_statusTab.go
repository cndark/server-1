package gamedata

var ConfStatusTab = &statusTabTable{}

type statusTab struct {
	Id   int32  `json:"id"`   // id
	Type string `json:"type"` // 类型
	Flag int32  `json:"flag"` // 参数1（类型）
}

type statusTabTable struct {
	items map[int32]*statusTab
}

func (self *statusTabTable) Load() {
	var arr []*statusTab
	if !load_json("statusTab.json", &arr) {
		return
	}

	items := make(map[int32]*statusTab)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *statusTabTable) Query(id int32) *statusTab {
	return self.items[id]
}

func (self *statusTabTable) Items() map[int32]*statusTab {
	return self.items
}
