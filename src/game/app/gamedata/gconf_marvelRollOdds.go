package gamedata

var ConfMarvelRollOdds = &marvelRollOddsTable{}

type marvelRollOdds struct {
	Id      int32  `json:"id"`      // id
	Group   string `json:"group"`   // 组
	BlankId int32  `json:"blankId"` // 格子
	Item    []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"item"` // 物品
	LevelRange []*struct {
		Low  int32 `json:"low"`
		High int32 `json:"high"`
	} `json:"levelRange"` // 等级区间
	RefOdds  int32 `json:"refOdds"`  // 刷新权重
	RollOdds int32 `json:"rollOdds"` // 许愿权重
	Upper    int32 `json:"upper"`    // 获得次数上限
}

type marvelRollOddsTable struct {
	items map[int32]*marvelRollOdds
}

func (self *marvelRollOddsTable) Load() {
	var arr []*marvelRollOdds
	if !load_json("marvelRollOdds.json", &arr) {
		return
	}

	items := make(map[int32]*marvelRollOdds)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *marvelRollOddsTable) Query(id int32) *marvelRollOdds {
	return self.items[id]
}

func (self *marvelRollOddsTable) Items() map[int32]*marvelRollOdds {
	return self.items
}
