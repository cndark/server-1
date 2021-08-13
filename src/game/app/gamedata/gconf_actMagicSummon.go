package gamedata

var ConfActMagicSummon = &actMagicSummonTable{}

type actMagicSummon struct {
	ConfGrp   int32 `json:"confGrp"` // 活动配置
	MagicHero []*struct {
		Hero   int32 `json:"hero"`
		Weight int32 `json:"weight"`
	} `json:"magicHero"` // 指定魔法英雄
	MagicExtraTimes int32  `json:"magicExtraTimes"` // 魔法保底次数
	DiamCnt         int32  `json:"diamCnt"`         // 本期可消耗钻石召唤次数
	DiamCost        int32  `json:"diamCost"`        // 召唤消耗钻石单价
	SummonGrp       string `json:"summonGrp"`       // 召唤组
	ExtraTimes      int32  `json:"extraTimes"`      // 常规保底次数
	ExtraGroup      string `json:"extraGroup"`      // 常规保底组
	Cost            []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 基础单价
	Discount     float32 `json:"discount"`     // 10连折扣
	CostFreeTime int32   `json:"costFreeTime"` // 免费时限(秒)
}

type actMagicSummonTable struct {
	items map[int32]*actMagicSummon
}

func (self *actMagicSummonTable) Load() {
	var arr []*actMagicSummon
	if !load_json("actMagicSummon.json", &arr) {
		return
	}

	items := make(map[int32]*actMagicSummon)

	for _, v := range arr {
		items[v.ConfGrp] = v
	}

	self.items = items
}

func (self *actMagicSummonTable) Query(confGrp int32) *actMagicSummon {
	return self.items[confGrp]
}

func (self *actMagicSummonTable) Items() map[int32]*actMagicSummon {
	return self.items
}
