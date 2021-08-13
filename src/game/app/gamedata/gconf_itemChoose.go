package gamedata

var ConfItemChoose = &itemChooseTable{}

type itemChoose struct {
	Id     int32 `json:"id"` // 任选id
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励库
}

type itemChooseTable struct {
	items map[int32]*itemChoose
}

func (self *itemChooseTable) Load() {
	var arr []*itemChoose
	if !load_json("itemChoose.json", &arr) {
		return
	}

	items := make(map[int32]*itemChoose)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *itemChooseTable) Query(id int32) *itemChoose {
	return self.items[id]
}

func (self *itemChooseTable) Items() map[int32]*itemChoose {
	return self.items
}
