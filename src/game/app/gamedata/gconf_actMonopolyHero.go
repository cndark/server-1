package gamedata

var ConfActMonopolyHero = &actMonopolyHeroTable{}

type actMonopolyHero struct {
	Seq     int32   `json:"seq"`     // 序列
	Group   int32   `json:"group"`   // 奖励组
	Monster []int32 `json:"monster"` // 怪物配置
	BossPos int32   `json:"bossPos"` // boss站位
	Reward  []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	Weight int32 `json:"weight"` // 出现权重
}

type actMonopolyHeroTable struct {
	items map[int32]*actMonopolyHero
}

func (self *actMonopolyHeroTable) Load() {
	var arr []*actMonopolyHero
	if !load_json("actMonopolyHero.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyHero)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMonopolyHeroTable) Query(seq int32) *actMonopolyHero {
	return self.items[seq]
}

func (self *actMonopolyHeroTable) Items() map[int32]*actMonopolyHero {
	return self.items
}
