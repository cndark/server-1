package gamedata

var ConfTargetDays = &targetDaysTable{}

type targetDays struct {
	Seq       int32 `json:"seq"`  // 序列
	Day       int32 `json:"day"`  // 天数
	Type      int32 `json:"type"` // 类型
	AttainTab []*struct {
		AttainId int32   `json:"attainId"`
		P2       float64 `json:"p2"`
	} `json:"attainTab"` // 条件状态
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
}

type targetDaysTable struct {
	items map[int32]*targetDays
}

func (self *targetDaysTable) Load() {
	var arr []*targetDays
	if !load_json("targetDays.json", &arr) {
		return
	}

	items := make(map[int32]*targetDays)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *targetDaysTable) Query(seq int32) *targetDays {
	return self.items[seq]
}

func (self *targetDaysTable) Items() map[int32]*targetDays {
	return self.items
}
