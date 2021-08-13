package gamedata

var ConfActMonopolyLevel = &actMonopolyLevelTable{}

type actMonopolyLevel struct {
	Seq         int32 `json:"seq"` // 关卡序列
	Lv          int32 `json:"lv"`  // 关卡等级
	EventWeight []*struct {
		Tp     int32 `json:"tp"`
		Weight int32 `json:"weight"`
		Grp    int32 `json:"grp"`
	} `json:"eventWeight"` // 奇遇生成类型和对应权重、奇遇类型组
	CreateEvent int32 `json:"createEvent"` // 奇遇生成区间
	EventMax    int32 `json:"eventMax"`    // 单一类型最大生成数量
	BaseReward  []*struct {
		Id     int32 `json:"id"`
		N      int64 `json:"n"`
		Weight int32 `json:"weight"`
	} `json:"baseReward"` // 基础奖励(和骰子次数挂钩)
	PassReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"passReward"` // 过关奖励
	ConfGrp int32 `json:"confGrp"` // 活动配置
}

type actMonopolyLevelTable struct {
	items map[int32]*actMonopolyLevel
}

func (self *actMonopolyLevelTable) Load() {
	var arr []*actMonopolyLevel
	if !load_json("actMonopolyLevel.json", &arr) {
		return
	}

	items := make(map[int32]*actMonopolyLevel)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMonopolyLevelTable) Query(seq int32) *actMonopolyLevel {
	return self.items[seq]
}

func (self *actMonopolyLevelTable) Items() map[int32]*actMonopolyLevel {
	return self.items
}
