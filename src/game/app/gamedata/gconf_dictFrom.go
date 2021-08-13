package gamedata

var ConfDictFrom = &dictFromTable{}

type dictFrom struct {
	Id   string `json:"id"`   // id
	Name string `json:"name"` // 名称
}

type dictFromTable struct {
	items map[string]*dictFrom
}

func (self *dictFromTable) Load() {
	var arr []*dictFrom
	if !load_json("dictFrom.json", &arr) {
		return
	}

	items := make(map[string]*dictFrom)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *dictFromTable) Query(id string) *dictFrom {
	return self.items[id]
}

func (self *dictFromTable) Items() map[string]*dictFrom {
	return self.items
}
