package gamedata

var ConfWorldLevelDraw = &worldLevelDrawTable{}

type worldLevelDraw struct {
	Id         int32   `json:"id"`         // 任务id
	Lv         int32   `json:"lv"`         // 达到的关卡
	NormalDrop int32   `json:"normalDrop"` // 普通掉落
	SeniorDrop []int32 `json:"seniorDrop"` // 5星掉落
}

type worldLevelDrawTable struct {
	items map[int32]*worldLevelDraw
}

func (self *worldLevelDrawTable) Load() {
	var arr []*worldLevelDraw
	if !load_json("worldLevelDraw.json", &arr) {
		return
	}

	items := make(map[int32]*worldLevelDraw)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *worldLevelDrawTable) Query(id int32) *worldLevelDraw {
	return self.items[id]
}

func (self *worldLevelDrawTable) Items() map[int32]*worldLevelDraw {
	return self.items
}
