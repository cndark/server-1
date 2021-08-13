package gamedata

var ConfMarvelRoll = &marvelRollTable{}

type marvelRoll struct {
	Group      string `json:"group"`    // 转盘组
	ModuleId   int32  `json:"moduleId"` // 功能板块，openId
	OpenStatus []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"openStatus"` // 开启需要的状态
	Cost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 基础单价
	Discount    float32 `json:"discount"` // 10连折扣
	RefreshCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"refreshCost"` // 刷新消耗
	ExtReward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"extReward"` // 每次额外获得
	BlankNum int32 `json:"blankNum"` // 格子数
}

type marvelRollTable struct {
	items map[string]*marvelRoll
}

func (self *marvelRollTable) Load() {
	var arr []*marvelRoll
	if !load_json("marvelRoll.json", &arr) {
		return
	}

	items := make(map[string]*marvelRoll)

	for _, v := range arr {
		items[v.Group] = v
	}

	self.items = items
}

func (self *marvelRollTable) Query(group string) *marvelRoll {
	return self.items[group]
}

func (self *marvelRollTable) Items() map[string]*marvelRoll {
	return self.items
}
