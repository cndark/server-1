package gamedata

var ConfHeroStarUp = &heroStarUpTable{}

type heroStarUp struct {
	Star     int32 `json:"star"`  // 星级
	MaxLv    int32 `json:"maxLv"` // 等级上限
	StarCost []*struct {
		Tp   int32 `json:"tp"`
		Star int32 `json:"star"`
		N    int32 `json:"n"`
	} `json:"starCost"` // 材料配方
	StarAddCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"starAddCost"` // 附加消耗
	StarRet []*struct {
		Type int32 `json:"type"`
		Id   int32 `json:"id"`
		N    int32 `json:"n"`
	} `json:"starRet"` // 升星消耗返还
	PropsRatio []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"propsRatio"` // 属性加成(替换)
	Sacrifice []*struct {
		Type int32 `json:"type"`
		Id   int32 `json:"id"`
		N    int32 `json:"n"`
	} `json:"sacrifice"` // 分解奖励(1道具|2自身|3通用碎片)
	StarInheritNum int32 `json:"starInheritNum"` // 星级继承碎片数量
}

type heroStarUpTable struct {
	items map[int32]*heroStarUp
}

func (self *heroStarUpTable) Load() {
	var arr []*heroStarUp
	if !load_json("heroStarUp.json", &arr) {
		return
	}

	items := make(map[int32]*heroStarUp)

	for _, v := range arr {
		items[v.Star] = v
	}

	self.items = items
}

func (self *heroStarUpTable) Query(star int32) *heroStarUp {
	return self.items[star]
}

func (self *heroStarUpTable) Items() map[int32]*heroStarUp {
	return self.items
}
