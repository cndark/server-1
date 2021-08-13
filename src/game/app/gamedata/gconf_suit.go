package gamedata

var ConfSuit = &suitTable{}

type suit struct {
	Id       int32 `json:"Id"` // 套装
	Property []*struct {
		N    int32   `json:"n"`
		Attr int32   `json:"attr"`
		Val  float32 `json:"val"`
	} `json:"property"` // 套装加成属性
}

type suitTable struct {
	items map[int32]*suit
}

func (self *suitTable) Load() {
	var arr []*suit
	if !load_json("suit.json", &arr) {
		return
	}

	items := make(map[int32]*suit)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *suitTable) Query(Id int32) *suit {
	return self.items[Id]
}

func (self *suitTable) Items() map[int32]*suit {
	return self.items
}
