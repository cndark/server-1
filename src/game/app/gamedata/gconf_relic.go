package gamedata

var ConfRelic = &relicTable{}

type relic struct {
	Id        int32  `json:"id"`       // id
	Txt_Name  string `json:"txt_Name"` // 名称
	Color     int32  `json:"color"`    // 品质颜色
	BaseProps []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"baseProps"` // 初始属性
	BasePropGrowth []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"basePropGrowth"` // 初始属性增长
	Active []*struct {
		Type int32 `json:"type"`
		N    int32 `json:"n"`
	} `json:"active"` // 激活条件-1阵营2职业3英雄
	HideProps []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"hideProps"` // 隐藏属性
	HidePropGrowth []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"hidePropGrowth"` // 隐藏属性增长
	Exp       int32 `json:"exp"`       // 经验需求
	SupplyExp int32 `json:"supplyExp"` // 白板提供经验
	StarLimit int32 `json:"starLimit"` // 星级上限
	Decompose []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"decompose"` // （品质>5）分解产出
}

type relicTable struct {
	items map[int32]*relic
}

func (self *relicTable) Load() {
	var arr []*relic
	if !load_json("relic.json", &arr) {
		return
	}

	items := make(map[int32]*relic)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *relicTable) Query(id int32) *relic {
	return self.items[id]
}

func (self *relicTable) Items() map[int32]*relic {
	return self.items
}
