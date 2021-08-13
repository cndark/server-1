package gamedata

var ConfAiBot = &aiBotTable{}

type aiBot struct {
	Seq  int32   `json:"seq"`  // 序列
	Type int32   `json:"type"` // 类型(1资源2英雄3建筑)
	Id   int32   `json:"id"`   // id
	N    float64 `json:"n"`    // n
}

type aiBotTable struct {
	items map[int32]*aiBot
}

func (self *aiBotTable) Load() {
	var arr []*aiBot
	if !load_json("aiBot.json", &arr) {
		return
	}

	items := make(map[int32]*aiBot)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *aiBotTable) Query(seq int32) *aiBot {
	return self.items[seq]
}

func (self *aiBotTable) Items() map[int32]*aiBot {
	return self.items
}
