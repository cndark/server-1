package gamedata

var ConfHeroChangeOdds = &heroChangeOddsTable{}

type heroChangeOdds struct {
	Id    int32 `json:"Id"`    // 序列
	Group int32 `json:"group"` // 组
	Hero  int32 `json:"hero"`  // 英雄
	Odds  int32 `json:"odds"`  // 权重
}

type heroChangeOddsTable struct {
	items map[int32]*heroChangeOdds
}

func (self *heroChangeOddsTable) Load() {
	var arr []*heroChangeOdds
	if !load_json("heroChangeOdds.json", &arr) {
		return
	}

	items := make(map[int32]*heroChangeOdds)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *heroChangeOddsTable) Query(Id int32) *heroChangeOdds {
	return self.items[Id]
}

func (self *heroChangeOddsTable) Items() map[int32]*heroChangeOdds {
	return self.items
}
