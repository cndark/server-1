package gamedata

var ConfGuildShrine = &guildShrineTable{}

type guildShrine struct {
	Id       int32  `json:"id"`       // 神灵id
	Txt_Name string `json:"txt_Name"` // 名字
	ElemId   int32  `json:"elemId"`   // 对应魂玉id
}

type guildShrineTable struct {
	items map[int32]*guildShrine
}

func (self *guildShrineTable) Load() {
	var arr []*guildShrine
	if !load_json("guildShrine.json", &arr) {
		return
	}

	items := make(map[int32]*guildShrine)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *guildShrineTable) Query(id int32) *guildShrine {
	return self.items[id]
}

func (self *guildShrineTable) Items() map[int32]*guildShrine {
	return self.items
}
