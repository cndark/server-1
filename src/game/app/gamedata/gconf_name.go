package gamedata

var ConfName = &nameTable{}

type name struct {
	NameId    int32  `json:"nameId"`    // 名字ID
	FirstName string `json:"firstName"` // 姓
	LastName  string `json:"lastName"`  // 名
	Language  string `json:"language"`  // 语言
}

type nameTable struct {
	items map[int32]*name
}

func (self *nameTable) Load() {
	var arr []*name
	if !load_json("name.json", &arr) {
		return
	}

	items := make(map[int32]*name)

	for _, v := range arr {
		items[v.NameId] = v
	}

	self.items = items
}

func (self *nameTable) Query(nameId int32) *name {
	return self.items[nameId]
}

func (self *nameTable) Items() map[int32]*name {
	return self.items
}
