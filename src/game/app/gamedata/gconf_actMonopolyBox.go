package gamedata

var ConfActMonopolyBox = &actMonopolyBoxTable{}

type actMonopolyBox struct {
	Seq       int32 `json:"seq"`   // 序列
	Group     int32 `json:"group"` // 奖励组
	BoxReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"boxReward"` // 宝箱奖励
	Weight int32 `json:"weight"` // 出现权重
}

type actMonopolyBoxTable struct {
	items map[int32]*actMonopolyBox
}

func (self *actMonopolyBoxTable) Load() {
	var arr []*actMonopolyBox
	if !load_json("actMonopolyBox.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyBox)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMonopolyBoxTable) Query(seq int32) *actMonopolyBox {
	return self.items[seq]
}

func (self *actMonopolyBoxTable) Items() map[int32]*actMonopolyBox {
	return self.items
}
