package gamedata

var ConfActSummonOdds = &actSummonOddsTable{}

type actSummonOdds struct {
	Seq       int32  `json:"seq"`       // 序列
	SummonGrp string `json:"summonGrp"` // 召唤组
	Item      []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"item"` // 道具
	Weight      int32 `json:"weight"`      // 权重
	GetCntLimit int32 `json:"getCntLimit"` // 抽取次数条件
	AddCntLimit int32 `json:"addCntLimit"` // 出货增加条件次数
}

type actSummonOddsTable struct {
	items map[int32]*actSummonOdds
}

func (self *actSummonOddsTable) Load() {
	var arr []*actSummonOdds
	if !load_json("actSummonOdds.json", &arr) {
		return
	}

	items := make(map[int32]*actSummonOdds)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actSummonOddsTable) Query(seq int32) *actSummonOdds {
	return self.items[seq]
}

func (self *actSummonOddsTable) Items() map[int32]*actSummonOdds {
	return self.items
}
