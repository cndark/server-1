package gamedata

var ConfCrossconf = &crossconfTable{}

type crossconf struct {
	Id   int32   `json:"id"`   // 分组id
	Svrs []int32 `json:"svrs"` // 服务器id
}

type crossconfTable struct {
	items map[int32]*crossconf
}

func (self *crossconfTable) Load() {
	var arr []*crossconf
	if !load_json("crossconf.json", &arr) {
		return
	}

	items := make(map[int32]*crossconf)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *crossconfTable) Query(id int32) *crossconf {
	return self.items[id]
}

func (self *crossconfTable) Items() map[int32]*crossconf {
	return self.items
}
