package gamedata

var ConfWorldLevel1000 = &worldLevel1000Table{}

type worldLevel1000 struct {
	Id        int32 `json:"id"`        // id
	RoundType int32 `json:"roundType"` // 回合类型
	Monster   []*struct {
		Id int32 `json:"id"`
		Lv int32 `json:"lv"`
	} `json:"monster"` // 关卡怪物
	PowerSwitch int32   `json:"powerSwitch"` // 强度等级阈值
	PowerRatio  float32 `json:"powerRatio"`  // 基准系数
	Reward      []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 通关奖励
	MinuteCurrency []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"minuteCurrency"` // 分钟挂机货币
	ExploreDrop int32 `json:"exploreDrop"` // 挂机掉落（挂机按分钟算）
	RewardDrop  int32 `json:"rewardDrop"`  // 通关掉落
}

type worldLevel1000Table struct {
	items map[int32]*worldLevel1000
}

func (self *worldLevel1000Table) Load() {
	var arr []*worldLevel1000
	if !load_json("worldLevel1000.json", &arr) {
		return
	}

	items := make(map[int32]*worldLevel1000)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *worldLevel1000Table) Query(id int32) *worldLevel1000 {
	return self.items[id]
}

func (self *worldLevel1000Table) Items() map[int32]*worldLevel1000 {
	return self.items
}
