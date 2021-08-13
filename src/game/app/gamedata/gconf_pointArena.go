package gamedata

var ConfPointArena = &pointArenaTable{}

type pointArena struct {
	Id        int32 `json:"id"` // 序号
	ArenaRank []*struct {
		Up   int32 `json:"up"`
		Down int32 `json:"down"`
	} `json:"arenaRank"` // pvp等级区间
	DailyMail int32 `json:"dailyMail"` // 每日结算邮件ID
	WeekMail  int32 `json:"weekMail"`  // 赛季结算邮件ID
}

type pointArenaTable struct {
	items map[int32]*pointArena
}

func (self *pointArenaTable) Load() {
	var arr []*pointArena
	if !load_json("pointArena.json", &arr) {
		return
	}

	items := make(map[int32]*pointArena)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *pointArenaTable) Query(id int32) *pointArena {
	return self.items[id]
}

func (self *pointArenaTable) Items() map[int32]*pointArena {
	return self.items
}
