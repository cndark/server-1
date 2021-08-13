package gamedata

var ConfHeroAwake = &heroAwakeTable{}

type heroAwake struct {
	Lv   int32 `json:"lv"` // 觉醒等级
	Attr []*struct {
		Elem int32   `json:"elem"`
		Id   int32   `json:"id"`
		Val  float64 `json:"val"`
	} `json:"attr"` // 觉醒属性
	Cost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"cost"` // 升级消耗
	LvCond  int32 `json:"lvCond"` // 鬼神等级要求
	CcyfRet []*struct {
		Id int32   `json:"id"`
		N  float64 `json:"n"`
	} `json:"ccyfRet"` // 英雄消耗返还
}

type heroAwakeTable struct {
	items map[int32]*heroAwake
}

func (self *heroAwakeTable) Load() {
	var arr []*heroAwake
	if !load_json("heroAwake.json", &arr) {
		return
	}

	items := make(map[int32]*heroAwake)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *heroAwakeTable) Query(lv int32) *heroAwake {
	return self.items[lv]
}

func (self *heroAwakeTable) Items() map[int32]*heroAwake {
	return self.items
}
