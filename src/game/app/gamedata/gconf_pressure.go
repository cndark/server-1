package gamedata

var ConfPressure = &pressureTable{}

type pressure struct {
	Id int32   `json:"id"` // id
	N  float64 `json:"n"`  // n
}

type pressureTable struct {
	items map[int32]*pressure
}

func (self *pressureTable) Load() {
	var arr []*pressure
	if !load_json("pressure.json", &arr) {
		return
	}

	items := make(map[int32]*pressure)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *pressureTable) Query(id int32) *pressure {
	return self.items[id]
}

func (self *pressureTable) Items() map[int32]*pressure {
	return self.items
}
