package gamedata

var ConfGlobalCross = &globalCrossTable{}

type globalCross struct {
	StructureNum int32 `json:"structureNum"` // 结构序列
}

type globalCrossTable struct {
	items map[int32]*globalCross
}

func (self *globalCrossTable) Load() {
	var arr []*globalCross
	if !load_json("globalCross.json", &arr) {
		return
	}

	items := make(map[int32]*globalCross)

	for _, v := range arr {
		items[v.StructureNum] = v
	}

	self.items = items
}

func (self *globalCrossTable) Query(structureNum int32) *globalCross {
	return self.items[structureNum]
}

func (self *globalCrossTable) Items() map[int32]*globalCross {
	return self.items
}
