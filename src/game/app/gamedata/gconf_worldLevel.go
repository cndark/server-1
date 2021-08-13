package gamedata

var ConfWorldLevel = &worldLevelTable{}

type worldLevel struct {
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

type worldLevelTable struct {
	items map[int32]*worldLevel
}

func (self *worldLevelTable) Load() {
	var arr []*worldLevel
	if !load_json("worldLevel.json", &arr) {
		return
	}

	items := make(map[int32]*worldLevel)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *worldLevelTable) Query(id int32) *worldLevel {
	return self.items[id]
}

func (self *worldLevelTable) Items() map[int32]*worldLevel {
	return self.items
}
