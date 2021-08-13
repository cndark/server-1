package gamedata

var ConfCalendar = &calendarTable{}

type calendar struct {
	Seq       int32  `json:"seq"`     // 序列
	Type      string `json:"type"`    // 类型
	ModName   string `json:"modName"` // 模块名
	StageList []*struct {
		Ts    string `json:"ts"`
		Stage string `json:"stage"`
	} `json:"stageList"` // 阶段列表
}

type calendarTable struct {
	items map[int32]*calendar
}

func (self *calendarTable) Load() {
	var arr []*calendar
	if !load_json("calendar.json", &arr) {
		return
	}

	items := make(map[int32]*calendar)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *calendarTable) Query(seq int32) *calendar {
	return self.items[seq]
}

func (self *calendarTable) Items() map[int32]*calendar {
	return self.items
}
