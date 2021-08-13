package gamedata

var ConfGlobalBattle = &globalBattleTable{}

type globalBattle struct {
	StructureNum    int32   `json:"structureNum"`   // 结构序列
	BaseDamRatio    float32 `json:"baseDamRatio"`   // 保底基础伤害
	BlockDamRatio   float32 `json:"blockDamRatio"`  // 格挡减伤
	CritRatioLimit  float32 `json:"critRatioLimit"` // 实际暴击率上限
	BattleLifeRatio []*struct {
		Stat int32   `json:"stat"`
		N    float32 `json:"n"`
	} `json:"battleLifeRatio"` // pvp生命加成
	EnergyLimit      int32 `json:"energyLimit"`    // 能量上限
	InitailEnergy    int32 `json:"initailEnergy"`  // 初始能量
	NormalToEnergy   int32 `json:"normalToEnergy"` // 普攻增加能量
	DefendToEnergy   int32 `json:"defendToEnergy"` // 受击增加能量
	MonsterBaseProps []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"monsterBaseProps"` // 怪物初始属性
}

type globalBattleTable struct {
	items map[int32]*globalBattle
}

func (self *globalBattleTable) Load() {
	var arr []*globalBattle
	if !load_json("globalBattle.json", &arr) {
		return
	}

	items := make(map[int32]*globalBattle)

	for _, v := range arr {
		items[v.StructureNum] = v
	}

	self.items = items
}

func (self *globalBattleTable) Query(structureNum int32) *globalBattle {
	return self.items[structureNum]
}

func (self *globalBattleTable) Items() map[int32]*globalBattle {
	return self.items
}
