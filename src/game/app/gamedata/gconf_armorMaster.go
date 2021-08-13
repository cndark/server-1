package gamedata

var ConfArmorMaster = &armorMasterTable{}

type armorMaster struct {
	Lv   int32 `json:"lv"` // 等级
	Prop []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"prop"` // 装备大师属性
}

type armorMasterTable struct {
	items map[int32]*armorMaster
}

func (self *armorMasterTable) Load() {
	var arr []*armorMaster
	if !load_json("armorMaster.json", &arr) {
		return
	}

	items := make(map[int32]*armorMaster)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *armorMasterTable) Query(lv int32) *armorMaster {
	return self.items[lv]
}

func (self *armorMasterTable) Items() map[int32]*armorMaster {
	return self.items
}
