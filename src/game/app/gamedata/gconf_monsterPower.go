package gamedata

var ConfMonsterPower = &monsterPowerTable{}

type monsterPower struct {
	Level     int32 `json:"level"` // 等级
	Star      int32 `json:"star"`  // 星级
	BaseProps []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"baseProps"` // 初始属性
}

type monsterPowerTable struct {
	items map[int32]*monsterPower
}

func (self *monsterPowerTable) Load() {
	var arr []*monsterPower
	if !load_json("monsterPower.json", &arr) {
		return
	}

	items := make(map[int32]*monsterPower)

	for _, v := range arr {
		items[v.Level] = v
	}

	self.items = items
}

func (self *monsterPowerTable) Query(level int32) *monsterPower {
	return self.items[level]
}

func (self *monsterPowerTable) Items() map[int32]*monsterPower {
	return self.items
}
