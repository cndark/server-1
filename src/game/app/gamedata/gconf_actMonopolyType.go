package gamedata

var ConfActMonopolyType = &actMonopolyTypeTable{}

type actMonopolyType struct {
	Tp   int32 `json:"tp"`   // 奇遇类型
	Time int32 `json:"time"` // 存在时间（秒）
}

type actMonopolyTypeTable struct {
	items map[int32]*actMonopolyType
}

func (self *actMonopolyTypeTable) Load() {
	var arr []*actMonopolyType
	if !load_json("actMonopolyType.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyType)

	for _, v := range arr {
		items[v.Tp] = v
	}

	self.items = items
}

func (self *actMonopolyTypeTable) Query(tp int32) *actMonopolyType {
	return self.items[tp]
}

func (self *actMonopolyTypeTable) Items() map[int32]*actMonopolyType {
	return self.items
}
