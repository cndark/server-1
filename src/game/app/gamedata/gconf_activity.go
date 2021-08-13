package gamedata

var ConfActivity = &activityTable{}

type activity struct {
	Name       string `json:"name"` // 活动名en
	ConfCycles []*struct {
		Cnt  int32  `json:"cnt"`
		Grps string `json:"grps"`
	} `json:"confCycles"` // 活动配置组: 循环次数~配置id列表| ….
	NeedDays     int32 `json:"needDays"`     // 需要的开服天数
	InitialPower int32 `json:"initialPower"` // 初始体力
}

type activityTable struct {
	items map[string]*activity
}

func (self *activityTable) Load() {
	var arr []*activity
	if !load_json("activity.json", &arr) {
		return
	}

	items := make(map[string]*activity)

	for _, v := range arr {
		items[v.Name] = v
	}

	self.items = items
}

func (self *activityTable) Query(name string) *activity {
	return self.items[name]
}

func (self *activityTable) Items() map[string]*activity {
	return self.items
}
