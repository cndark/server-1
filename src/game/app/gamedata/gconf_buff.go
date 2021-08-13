package gamedata

var ConfBuff = &buffTable{}

type buff struct {
	Id int32 `json:"id"` // ID
}

type buffTable struct {
	items map[int32]*buff
}

func (self *buffTable) Load() {
	var arr []*buff
	if !load_json("buff.json", &arr) {
		return
	}

	items := make(map[int32]*buff)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *buffTable) Query(id int32) *buff {
	return self.items[id]
}

func (self *buffTable) Items() map[int32]*buff {
	return self.items
}
