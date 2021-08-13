package gamedata

var ConfWorldBoss = &worldBossTable{}

type worldBoss struct {
	BossId       int32 `json:"bossId"`    // BossId
	MonsterId    int32 `json:"monsterId"` // monsterId
	Lv           int32 `json:"lv"`        // boss等级
	RoundType    int32 `json:"roundType"` // 回合类型
	DamageReward []*struct {
		N  float64 `json:"n"`
		Id int32   `json:"id"`
	} `json:"damageReward"` // 个人阶段伤害奖励
	MaxDamageReward []*struct {
		N  float64 `json:"n"`
		Id int32   `json:"id"`
	} `json:"maxDamageReward"` // 最高伤害奖励
}

type worldBossTable struct {
	items map[int32]*worldBoss
}

func (self *worldBossTable) Load() {
	var arr []*worldBoss
	if !load_json("worldBoss.json", &arr) {
		return
	}

	items := make(map[int32]*worldBoss)

	for _, v := range arr {
		items[v.BossId] = v
	}

	self.items = items
}

func (self *worldBossTable) Query(bossId int32) *worldBoss {
	return self.items[bossId]
}

func (self *worldBossTable) Items() map[int32]*worldBoss {
	return self.items
}
