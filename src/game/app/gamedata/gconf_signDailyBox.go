package gamedata

var ConfSignDailyBox = &signDailyBoxTable{}

type signDailyBox struct {
	Day    int32 `json:"day"` // 天数
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type signDailyBoxTable struct {
	items map[int32]*signDailyBox
}

func (self *signDailyBoxTable) Load() {
	var arr []*signDailyBox
	if !load_json("signDailyBox.json", &arr) {
		return
	}

	items := make(map[int32]*signDailyBox)

	for _, v := range arr {
		items[v.Day] = v
	}

	self.items = items
}

func (self *signDailyBoxTable) Query(day int32) *signDailyBox {
	return self.items[day]
}

func (self *signDailyBoxTable) Items() map[int32]*signDailyBox {
	return self.items
}
