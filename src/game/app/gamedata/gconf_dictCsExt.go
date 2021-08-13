package gamedata

var ConfDictCsExt = &dictCsExtTable{}

type dictCsExt struct {
	Id       string `json:"id"`       // id
	Txt_Name string `json:"txt_Name"` // 名称
}

type dictCsExtTable struct {
	items map[string]*dictCsExt
}

func (self *dictCsExtTable) Load() {
	var arr []*dictCsExt
	if !load_json("dictCsExt.json", &arr) {
		return
	}

	items := make(map[string]*dictCsExt)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *dictCsExtTable) Query(id string) *dictCsExt {
	return self.items[id]
}

func (self *dictCsExtTable) Items() map[string]*dictCsExt {
	return self.items
}
