package gamedata

var ConfGuildHarbor = &guildHarborTable{}

type guildHarbor struct {
	Level      int32   `json:"level"`      // 港口等级
	LevelExp   int64   `json:"levelExp"`   // 升级经验
	OutAdd     float32 `json:"outAdd"`     // 产出加成（百分比）
	OrderLimit int32   `json:"orderLimit"` // 订单上限
	TimeReduce float32 `json:"timeReduce"` // 消耗时间减少(百分比)
}

type guildHarborTable struct {
	items map[int32]*guildHarbor
}

func (self *guildHarborTable) Load() {
	var arr []*guildHarbor
	if !load_json("guildHarbor.json", &arr) {
		return
	}

	items := make(map[int32]*guildHarbor)

	for _, v := range arr {
		items[v.Level] = v
	}

	self.items = items
}

func (self *guildHarborTable) Query(level int32) *guildHarbor {
	return self.items[level]
}

func (self *guildHarborTable) Items() map[int32]*guildHarbor {
	return self.items
}
