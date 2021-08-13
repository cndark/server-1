package gamedata

var ConfActRushLocal = &actRushLocalTable{}

type actRushLocal struct {
	Seq       int32 `json:"seq"`       // seq
	RankId    int32 `json:"rankId"`    // 榜单id
	RewardGrp int32 `json:"rewardGrp"` // 奖励组
	ConfGrp   int32 `json:"confGrp"`   // 活动配置
}

type actRushLocalTable struct {
	items map[int32]*actRushLocal
}

func (self *actRushLocalTable) Load() {
	var arr []*actRushLocal
	if !load_json("actRushLocal.json", &arr) {
		return
	}

	items := make(map[int32]*actRushLocal)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actRushLocalTable) Query(seq int32) *actRushLocal {
	return self.items[seq]
}

func (self *actRushLocalTable) Items() map[int32]*actRushLocal {
	return self.items
}
