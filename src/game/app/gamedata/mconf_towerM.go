package gamedata

var ConfTowerM = &towerTableM{}

type towerTableM struct {
	items map[int32]*tower
}

func (self *towerTableM) Load() {
	items := make(map[int32]*tower)

	for _, v := range ConfTower.items {
		items[v.Num] = v
	}

	for _, v := range ConfTower500.items {
		items[v.Num] = (*tower)(v)
	}

	self.items = items
}

func (self *towerTableM) Query(id int32) *tower {
	return self.items[id]
}

func (self *towerTableM) Items() map[int32]*tower {
	return self.items
}
