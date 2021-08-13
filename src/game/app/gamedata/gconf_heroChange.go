package gamedata

var ConfHeroChange = &heroChangeTable{}

type heroChange struct {
	Hero  int32 `json:"hero"` // 英雄
	Group []*struct {
		Star int32 `json:"star"`
		Grp  int32 `json:"grp"`
	} `json:"group"` // 转换组
}

type heroChangeTable struct {
	items map[int32]*heroChange
}

func (self *heroChangeTable) Load() {
	var arr []*heroChange
	if !load_json("heroChange.json", &arr) {
		return
	}

	items := make(map[int32]*heroChange)

	for _, v := range arr {
		items[v.Hero] = v
	}

	self.items = items
}

func (self *heroChangeTable) Query(hero int32) *heroChange {
	return self.items[hero]
}

func (self *heroChangeTable) Items() map[int32]*heroChange {
	return self.items
}
