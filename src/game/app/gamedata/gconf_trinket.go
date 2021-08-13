package gamedata

var ConfTrinket = &trinketTable{}

type trinket struct {
	Lv           int32  `json:"lv"`           // 等级
	Txt_Name     string `json:"txt_Name"`     // 名称
	Color        int32  `json:"color"`        // 品质颜色
	AttributeNum int32  `json:"attributeNum"` // 属性条数
	UpCost       []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"upCost"` // 升级消耗
	UpRet []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"upRet"` // 升级消耗返还
	LockCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"lockCost"` // 锁定额外消耗
	TransformCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"transformCost"` // 属性转换消耗
	Prop []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
		Wt  int32   `json:"wt"`
	} `json:"prop"` // 饰品属性
}

type trinketTable struct {
	items map[int32]*trinket
}

func (self *trinketTable) Load() {
	var arr []*trinket
	if !load_json("trinket.json", &arr) {
		return
	}

	items := make(map[int32]*trinket)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *trinketTable) Query(lv int32) *trinket {
	return self.items[lv]
}

func (self *trinketTable) Items() map[int32]*trinket {
	return self.items
}
