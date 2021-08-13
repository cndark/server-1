package gamedata

var ConfActMazeLv = &actMazeLvTable{}

type actMazeLv struct {
	Seq            int32   `json:"seq"`            // 序列
	Lv             int32   `json:"lv"`             // 关卡id
	RandEventGrp   []int32 `json:"randEventGrp"`   // 随机事件组
	MonsterLvRatio float32 `json:"monsterLvRatio"` // 怪物等级系数
	PowerSwitch    int32   `json:"powerSwitch"`    // 强度等级阈值
	PowerRatio     float32 `json:"powerRatio"`     // 基准系数
	ExitPos        []int32 `json:"exitPos"`        // 出口位置随机
	ConfGrp        int32   `json:"confGrp"`        // 活动组
}

type actMazeLvTable struct {
	items map[int32]*actMazeLv
}

func (self *actMazeLvTable) Load() {
	var arr []*actMazeLv
	if !load_json("actMazeLv.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeLv)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMazeLvTable) Query(seq int32) *actMazeLv {
	return self.items[seq]
}

func (self *actMazeLvTable) Items() map[int32]*actMazeLv {
	return self.items
}
