package gamedata

var ConfDrop = &dropTable{}

type drop struct {
	DropId    int32 `json:"dropId"` // 掉落ID
	FixedDrop []*struct {
		Id int32   `json:"id"`
		N  float64 `json:"n"`
	} `json:"fixedDrop"` // 固定掉落
	RollDrop []*struct {
		Grp  int32 `json:"grp"`
		N    int32 `json:"n"`
		Prob int32 `json:"prob"`
	} `json:"rollDrop"` // 随机物品组
	CondDrop []int32 `json:"condDrop"` // 条件掉落
}

type dropTable struct {
	items map[int32]*drop
}

func (self *dropTable) Load() {
	var arr []*drop
	if !load_json("drop.json", &arr) {
		return
	}

	items := make(map[int32]*drop)

	for _, v := range arr {
		items[v.DropId] = v
	}

	self.items = items
}

func (self *dropTable) Query(dropId int32) *drop {
	return self.items[dropId]
}

func (self *dropTable) Items() map[int32]*drop {
	return self.items
}
