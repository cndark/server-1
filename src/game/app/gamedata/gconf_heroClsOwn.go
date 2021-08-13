package gamedata

var ConfHeroClsOwn = &heroClsOwnTable{}

type heroClsOwn struct {
	Lv   int32 `json:"lv"` // 突破等级
	Cost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"cost"` // 常规消耗
	CostOwn int32 `json:"costOwn"` // 专属消耗
	LvAttr  []*struct {
		Elem int32   `json:"elem"`
		Id   int32   `json:"id"`
		Val  float64 `json:"val"`
	} `json:"lvAttr"` // 每级提升，读当前
	ExtAttr []*struct {
		Elem int32   `json:"elem"`
		Id   int32   `json:"id"`
		Val  float64 `json:"val"`
	} `json:"extAttr"` // 每多少级额外提升的属性
	CcyfRet []*struct {
		Id int32   `json:"id"`
		N  float64 `json:"n"`
	} `json:"ccyfRet"` // 英雄消耗返还
}

type heroClsOwnTable struct {
	items map[int32]*heroClsOwn
}

func (self *heroClsOwnTable) Load() {
	var arr []*heroClsOwn
	if !load_json("heroClsOwn.json", &arr) {
		return
	}

	items := make(map[int32]*heroClsOwn)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *heroClsOwnTable) Query(lv int32) *heroClsOwn {
	return self.items[lv]
}

func (self *heroClsOwnTable) Items() map[int32]*heroClsOwn {
	return self.items
}
