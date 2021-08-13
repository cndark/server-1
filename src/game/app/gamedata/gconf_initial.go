package gamedata

var ConfInitial = &initialTable{}

type initial struct {
	Num                int32 `json:"num"`                // 序列
	InitialPlayerLevel int32 `json:"initialPlayerLevel"` // 初始玩家等级
	InitialCurrency    []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"initialCurrency"` // 初始货币
	InitialHero  []int32 `json:"initialHero"`  // 初始英雄
	InitialBeast []int32 `json:"initialBeast"` // 初始妖兽
	InitialItem  []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"initialItem"` // 初始道具与数量
	InitialModule []int32 `json:"initialModule"` // 初始模块
	InitialMailId int32   `json:"initialMailId"` // 初始邮件
}

type initialTable struct {
	items map[int32]*initial
}

func (self *initialTable) Load() {
	var arr []*initial
	if !load_json("initial.json", &arr) {
		return
	}

	items := make(map[int32]*initial)

	for _, v := range arr {
		items[v.Num] = v
	}

	self.items = items
}

func (self *initialTable) Query(num int32) *initial {
	return self.items[num]
}

func (self *initialTable) Items() map[int32]*initial {
	return self.items
}
