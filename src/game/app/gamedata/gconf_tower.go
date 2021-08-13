package gamedata

var ConfTower = &towerTable{}

type tower struct {
	Num       int32 `json:"num"`       // 层数
	RoundType int32 `json:"roundType"` // 回合类型
	Monster   []*struct {
		Id int32 `json:"id"`
		Lv int32 `json:"lv"`
	} `json:"monster"` // 守关怪
	PowerSwitch int32   `json:"powerSwitch"` // 强度等级阈值
	PowerRatio  float32 `json:"powerRatio"`  // 基准系数
	Reward      []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 通关奖励
	RaidReward []*struct {
		Id  int32   `json:"id"`
		N   int32   `json:"n"`
		Odd float32 `json:"odd"`
	} `json:"raidReward"` // 扫荡奖励
}

type towerTable struct {
	items map[int32]*tower
}

func (self *towerTable) Load() {
	var arr []*tower
	if !load_json("tower.json", &arr) {
		return
	}

	items := make(map[int32]*tower)

	for _, v := range arr {
		items[v.Num] = v
	}

	self.items = items
}

func (self *towerTable) Query(num int32) *tower {
	return self.items[num]
}

func (self *towerTable) Items() map[int32]*tower {
	return self.items
}
