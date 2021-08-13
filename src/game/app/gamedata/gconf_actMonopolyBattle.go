package gamedata

var ConfActMonopolyBattle = &actMonopolyBattleTable{}

type actMonopolyBattle struct {
	Seq     int32   `json:"seq"`     // 序列
	Group   int32   `json:"group"`   // 奖励组
	Monster []int32 `json:"monster"` // 怪物配置
	Reward  []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	Weight int32 `json:"weight"` // 出现权重
}

type actMonopolyBattleTable struct {
	items map[int32]*actMonopolyBattle
}

func (self *actMonopolyBattleTable) Load() {
	var arr []*actMonopolyBattle
	if !load_json("actMonopolyBattle.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyBattle)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMonopolyBattleTable) Query(seq int32) *actMonopolyBattle {
	return self.items[seq]
}

func (self *actMonopolyBattleTable) Items() map[int32]*actMonopolyBattle {
	return self.items
}
