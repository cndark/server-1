package gamedata

var ConfWorldLevelM = &worldLevelTableM{}

type worldLevelTableM struct {
	items map[int32]*worldLevel
}

func (self *worldLevelTableM) Load() {
	items := make(map[int32]*worldLevel)

	for _, v := range ConfWorldLevel.items {
		items[v.Id] = v
	}

	for _, v := range ConfWorldLevel500.items {
		items[v.Id] = (*worldLevel)(v)
	}

	for _, v := range ConfWorldLevel1000.items {
		items[v.Id] = (*worldLevel)(v)
	}

	self.items = items
}

func (self *worldLevelTableM) Query(id int32) *worldLevel {
	return self.items[id]
}

func (self *worldLevelTableM) Items() map[int32]*worldLevel {
	return self.items
}
