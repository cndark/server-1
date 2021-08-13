package gamedata

var ConfGuildTech = &guildTechTable{}

type guildTech struct {
	Id       int32 `json:"id"`       // 科技id
	Job      int32 `json:"job"`      // 职业
	LevelMax int32 `json:"levelMax"` // 等级上限
	InitProp []*struct {
		Id   int32   `json:"id"`
		Val  float32 `json:"val"`
		Grow float32 `json:"grow"`
	} `json:"initProp"` // 初始属性
	PreId []*struct {
		Id    int32 `json:"id"`
		Level int32 `json:"level"`
	} `json:"preId"` // 前置科技id
	UpCost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"upCost"` // 科技初始消耗
	UpCostGrowth []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"upCostGrowth"` // 消耗增长
}

type guildTechTable struct {
	items map[int32]*guildTech
}

func (self *guildTechTable) Load() {
	var arr []*guildTech
	if !load_json("guildTech.json", &arr) {
		return
	}

	items := make(map[int32]*guildTech)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *guildTechTable) Query(id int32) *guildTech {
	return self.items[id]
}

func (self *guildTechTable) Items() map[int32]*guildTech {
	return self.items
}
