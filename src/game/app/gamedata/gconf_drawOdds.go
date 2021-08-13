package gamedata

var ConfDrawOdds = &drawOddsTable{}

type drawOdds struct {
	Id        int32  `json:"Id"`        // 序列
	DrawGroup string `json:"drawGroup"` // 抽奖组
	Item      []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"item"` // 道具id
	BasicOdds   int32 `json:"basicOdds"` // 基础权重
	DrawCntOdds []*struct {
		Cnt int32 `json:"cnt"`
		Odd int32 `json:"odd"`
	} `json:"drawCntOdds"` // 动态权重（召唤次数）
	LevelOdds []*struct {
		Lv  int32 `json:"lv"`
		Odd int32 `json:"odd"`
	} `json:"levelOdds"` // 动态权重（等级）
	PayOdds []*struct {
		Diam int32 `json:"diam"`
		Odd  int32 `json:"odd"`
	} `json:"payOdds"` // 动态权重（付费）
	IsDelivery int32 `json:"isDelivery"` // 暂不投放
}

type drawOddsTable struct {
	items map[int32]*drawOdds
}

func (self *drawOddsTable) Load() {
	var arr []*drawOdds
	if !load_json("drawOdds.json", &arr) {
		return
	}

	items := make(map[int32]*drawOdds)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *drawOddsTable) Query(Id int32) *drawOdds {
	return self.items[Id]
}

func (self *drawOddsTable) Items() map[int32]*drawOdds {
	return self.items
}
