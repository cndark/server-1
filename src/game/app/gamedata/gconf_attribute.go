package gamedata

var ConfAttribute = &attributeTable{}

type attribute struct {
	AttributeId   int32   `json:"attributeId"`   // 属性ID
	PowerType     int32   `json:"powerType"`     // 战力系数类型
	HeroPower     float32 `json:"heroPower"`     // 英雄战力系数
	EquipPower    float32 `json:"equipPower"`    // 装备战力系数
	SoulPower     float32 `json:"soulPower"`     // 护灵战力系数
	ArtifactPower float32 `json:"artifactPower"` // 神器战力系数
	GemPower      float32 `json:"gemPower"`      // 宝石战力系数
}

type attributeTable struct {
	items map[int32]*attribute
}

func (self *attributeTable) Load() {
	var arr []*attribute
	if !load_json("attribute.json", &arr) {
		return
	}

	items := make(map[int32]*attribute)

	for _, v := range arr {
		items[v.AttributeId] = v
	}

	self.items = items
}

func (self *attributeTable) Query(attributeId int32) *attribute {
	return self.items[attributeId]
}

func (self *attributeTable) Items() map[int32]*attribute {
	return self.items
}
