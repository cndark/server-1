package gamedata

var ConfWarCupGuessTaskAttain = &warCupGuessTaskAttainTable{}

type warCupGuessTaskAttain struct {
	Id   int32 `json:"id"`   // id
	Cond int32 `json:"cond"` // 条件id
	P1   int32 `json:"p1"`   // p1
}

type warCupGuessTaskAttainTable struct {
	items map[int32]*warCupGuessTaskAttain
}

func (self *warCupGuessTaskAttainTable) Load() {
	var arr []*warCupGuessTaskAttain
	if !load_json("warCupGuessTaskAttain.json", &arr) {
		return
	}

	items := make(map[int32]*warCupGuessTaskAttain)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *warCupGuessTaskAttainTable) Query(id int32) *warCupGuessTaskAttain {
	return self.items[id]
}

func (self *warCupGuessTaskAttainTable) Items() map[int32]*warCupGuessTaskAttain {
	return self.items
}
