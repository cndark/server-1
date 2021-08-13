package gamedata

var ConfHeroUp = &heroUpTable{}

type heroUp struct {
	Lv   int32 `json:"lv"` // 等级
	Cost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"cost"` // 英雄升级消耗
	Ret []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"ret"` // 英雄消耗返还
}

type heroUpTable struct {
	items map[int32]*heroUp
}

func (self *heroUpTable) Load() {
	var arr []*heroUp
	if !load_json("heroUp.json", &arr) {
		return
	}

	items := make(map[int32]*heroUp)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *heroUpTable) Query(lv int32) *heroUp {
	return self.items[lv]
}

func (self *heroUpTable) Items() map[int32]*heroUp {
	return self.items
}
