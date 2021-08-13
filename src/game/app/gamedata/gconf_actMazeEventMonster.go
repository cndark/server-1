package gamedata

var ConfActMazeEventMonster = &actMazeEventMonsterTable{}

type actMazeEventMonster struct {
	Seq     int32   `json:"seq"`     // 序列
	Type    int32   `json:"type"`    // 怪物类型，0小怪，1boss
	Monster []int32 `json:"monster"` // 怪物队伍
	BossPos int32   `json:"bossPos"` // boss位置
	BuffNum int32   `json:"buffNum"` // 掉落buff个数
	Reward  []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 胜利奖励
	ShowLvCond int32 `json:"showLvCond"` // 多少层才出现
	Weight     int32 `json:"weight"`     // 出现权重
}

type actMazeEventMonsterTable struct {
	items map[int32]*actMazeEventMonster
}

func (self *actMazeEventMonsterTable) Load() {
	var arr []*actMazeEventMonster
	if !load_json("actMazeEventMonster.json", &arr) {
		return
	}

	items := make(map[int32]*actMazeEventMonster)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actMazeEventMonsterTable) Query(seq int32) *actMazeEventMonster {
	return self.items[seq]
}

func (self *actMazeEventMonsterTable) Items() map[int32]*actMazeEventMonster {
	return self.items
}
