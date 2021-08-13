package gamedata

var ConfMonster = &monsterTable{}

type monster struct {
	Id        int32  `json:"id"`        // ID
	Name      string `json:"name"`      // 后台名称
	Star      int32  `json:"star"`      // 英雄初始星级
	StarLimit int32  `json:"starLimit"` // 英雄星级上限
	Fragment  []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"fragment"` // 英雄碎片
	Elem          int32 `json:"elem"`        // 阵营
	JobId         int32 `json:"jobId"`       // 职业
	PlayerHead    int32 `json:"playerHead"`  // 激活头像
	ModifyPower   int32 `json:"modifyPower"` // 战斗力修正
	HeroBaseProps []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"heroBaseProps"` // 英雄初始属性
	HeroPropGrowth []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"heroPropGrowth"` // 英雄每级增长属性
	MonsterPropsRatio []float32 `json:"monsterPropsRatio"` // 怪物属性系数
	ActSummonCnt      int32     `json:"actSummonCnt"`      // 主题召唤保底次数
}

type monsterTable struct {
	items map[int32]*monster
}

func (self *monsterTable) Load() {
	var arr []*monster
	if !load_json("monster.json", &arr) {
		return
	}

	items := make(map[int32]*monster)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *monsterTable) Query(id int32) *monster {
	return self.items[id]
}

func (self *monsterTable) Items() map[int32]*monster {
	return self.items
}
