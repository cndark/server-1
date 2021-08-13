package gamedata

var ConfActHeroSkin = &actHeroSkinTable{}

type actHeroSkin struct {
	ConfGrp int32 `json:"confGrp"` // 活动配置
}

type actHeroSkinTable struct {
	items map[int32]*actHeroSkin
}

func (self *actHeroSkinTable) Load() {
	var arr []*actHeroSkin
	if !load_json("actHeroSkin.json", &arr) {
		return
	}

	items := make(map[int32]*actHeroSkin)

	for _, v := range arr {
		items[v.ConfGrp] = v
	}

	self.items = items
}

func (self *actHeroSkinTable) Query(confGrp int32) *actHeroSkin {
	return self.items[confGrp]
}

func (self *actHeroSkinTable) Items() map[int32]*actHeroSkin {
	return self.items
}
