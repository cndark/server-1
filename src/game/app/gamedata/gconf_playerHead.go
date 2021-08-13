package gamedata

var ConfPlayerHead = &playerHeadTable{}

type playerHead struct {
	Id string `json:"id"` // 对应英雄id
}

type playerHeadTable struct {
	items map[string]*playerHead
}

func (self *playerHeadTable) Load() {
	var arr []*playerHead
	if !load_json("playerHead.json", &arr) {
		return
	}

	items := make(map[string]*playerHead)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *playerHeadTable) Query(id string) *playerHead {
	return self.items[id]
}

func (self *playerHeadTable) Items() map[string]*playerHead {
	return self.items
}
