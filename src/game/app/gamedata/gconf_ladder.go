package gamedata

var ConfLadder = &ladderTable{}

type ladder struct {
	Id        int32 `json:"id"` // 序号
	RankRange []*struct {
		Up   int32 `json:"up"`
		Down int32 `json:"down"`
	} `json:"rankRange"` // 排名区间
	Reward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
}

type ladderTable struct {
	items map[int32]*ladder
}

func (self *ladderTable) Load() {
	var arr []*ladder
	if !load_json("ladder.json", &arr) {
		return
	}

	items := make(map[int32]*ladder)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *ladderTable) Query(id int32) *ladder {
	return self.items[id]
}

func (self *ladderTable) Items() map[int32]*ladder {
	return self.items
}
