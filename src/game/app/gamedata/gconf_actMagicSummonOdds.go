package gamedata

var ConfActMagicSummonOdds = &actMagicSummonOddsTable{}

type actMagicSummonOdds struct {
	Seq            int32  `json:"seq"`            // 序列
	SummonGrp      string `json:"summonGrp"`      // 召唤组
	RewardId       int32  `json:"rewardId"`       // 奖励
	RewardNum      int32  `json:"rewardNum"`      // 奖励数量
	Weight         int32  `json:"weight"`         // 权重
	MagicTimesCond int32  `json:"magicTimesCond"` // 魔法召唤次数条件
}

type actMagicSummonOddsTable struct {
	items map[int32]*actMagicSummonOdds
}

func (self *actMagicSummonOddsTable) Load() {
	var arr []*actMagicSummonOdds
	if !load_json("actMagicSummonOdds.json", &arr) {
		return
	}

	items := make(map[int32]*actMagicSummonOdds)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMagicSummonOddsTable) Query(seq int32) *actMagicSummonOdds {
	return self.items[seq]
}

func (self *actMagicSummonOddsTable) Items() map[int32]*actMagicSummonOdds {
	return self.items
}
