package gamedata

var ConfSignDaily = &signDailyTable{}

type signDaily struct {
	Day     int32 `json:"day"` // 天数
	Rewards []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"rewards"` // 奖励
}

type signDailyTable struct {
	items map[int32]*signDaily
}

func (self *signDailyTable) Load() {
	var arr []*signDaily
	if !load_json("signDaily.json", &arr) {
		return
	}

	items := make(map[int32]*signDaily)

	for _, v := range arr {
		items[v.Day] = v
	}

	self.items = items
}

func (self *signDailyTable) Query(day int32) *signDaily {
	return self.items[day]
}

func (self *signDailyTable) Items() map[int32]*signDaily {
	return self.items
}
