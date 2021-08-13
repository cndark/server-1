package gamedata

var ConfActMazeRank = &actMazeRankTable{}

type actMazeRank struct {
	Id   int32 `json:"id"` // 序列
	Rank []*struct {
		Low  int32 `json:"low"`
		High int32 `json:"high"`
	} `json:"rank"` // 排名范围
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	ConfGrp int32 `json:"confGrp"` // 活动组
}

type actMazeRankTable struct {
	items map[int32]*actMazeRank
}

func (self *actMazeRankTable) Load() {
	var arr []*actMazeRank
	if !load_json("actMazeRank.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeRank)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actMazeRankTable) Query(id int32) *actMazeRank {
	return self.items[id]
}

func (self *actMazeRankTable) Items() map[int32]*actMazeRank {
	return self.items
}
