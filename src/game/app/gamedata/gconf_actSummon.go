package gamedata

var ConfActSummon = &actSummonTable{}

type actSummon struct {
	ConfGrp    int32 `json:"confGrp"` // 活动配置
	Type       int32 `json:"type"`    // 召唤类型
	DesireHero []*struct {
		Hero   int32 `json:"hero"`
		Weight int32 `json:"weight"`
	} `json:"desireHero"` // 心愿英雄选1个（也属于UP英雄）
	UpHero []*struct {
		Hero   int32 `json:"hero"`
		Weight int32 `json:"weight"`
	} `json:"upHero"` // UP英雄选2个
	UpTimes    int32  `json:"upTimes"`    // UP保底次数
	DiamCnt    int32  `json:"diamCnt"`    // 本期可消耗钻石召唤次数
	DiamCost   int32  `json:"diamCost"`   // 召唤消耗钻石单价
	SummonGrp  string `json:"summonGrp"`  // 召唤组
	ExtraTimes int32  `json:"extraTimes"` // 保底次数
	ExtraGroup string `json:"extraGroup"` // 保底组
	Cost       []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 基础单价
	Discount     float32 `json:"discount"`     // 10连折扣
	CostFreeTime int32   `json:"costFreeTime"` // 免费时限(秒)
}

type actSummonTable struct {
	items map[int32]*actSummon
}

func (self *actSummonTable) Load() {
	var arr []*actSummon
	if !load_json("actSummon.json", &arr) {
		return
	}

	items := make(map[int32]*actSummon)

	for _, v := range arr {
		items[v.ConfGrp] = v
	}

	self.items = items
}

func (self *actSummonTable) Query(confGrp int32) *actSummon {
	return self.items[confGrp]
}

func (self *actSummonTable) Items() map[int32]*actSummon {
	return self.items
}
