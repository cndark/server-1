package gamedata

var ConfOpen = &openTable{}

type open struct {
	ModuleId int32   `json:"moduleId"` // 模块id
	Cond     int32   `json:"cond"`     // 条件ID
	P1       int32   `json:"p1"`       // 参数1
	P2       float64 `json:"p2"`       // 参数2
	Reward   []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 开启奖励
}

type openTable struct {
	items map[int32]*open
}

func (self *openTable) Load() {
	var arr []*open
	if !load_json("open.json", &arr) {
		return
	}

	items := make(map[int32]*open)

	for _, v := range arr {
		items[v.ModuleId] = v
	}

	self.items = items
}

func (self *openTable) Query(moduleId int32) *open {
	return self.items[moduleId]
}

func (self *openTable) Items() map[int32]*open {
	return self.items
}
