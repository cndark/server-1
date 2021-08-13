package gamedata

var ConfDaySign = &daySignTable{}

type daySign struct {
	Day    int32 `json:"day"` // 天数
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type daySignTable struct {
	items map[int32]*daySign
}

func (self *daySignTable) Load() {
	var arr []*daySign
	if !load_json("daySign.json", &arr) {
		return
	}

	items := make(map[int32]*daySign)

	for _, v := range arr {
		items[v.Day] = v
	}

	self.items = items
}

func (self *daySignTable) Query(day int32) *daySign {
	return self.items[day]
}

func (self *daySignTable) Items() map[int32]*daySign {
	return self.items
}
