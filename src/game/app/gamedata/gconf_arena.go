package gamedata

var ConfArena = &arenaTable{}

type arena struct {
	Id        int32 `json:"id"` // 序号
	RankRange []*struct {
		Up   int32 `json:"up"`
		Down int32 `json:"down"`
	} `json:"rankRange"` // 排名区间
	DailyMail int32 `json:"dailyMail"` // 每日结算邮件ID
	WeekMail  int32 `json:"weekMail"`  // 赛季结算邮件ID
}

type arenaTable struct {
	items map[int32]*arena
}

func (self *arenaTable) Load() {
	var arr []*arena
	if !load_json("arena.json", &arr) {
		return
	}

	items := make(map[int32]*arena)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *arenaTable) Query(id int32) *arena {
	return self.items[id]
}

func (self *arenaTable) Items() map[int32]*arena {
	return self.items
}
