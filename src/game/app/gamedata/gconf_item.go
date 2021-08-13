package gamedata

var ConfItem = &itemTable{}

type item struct {
	Id            int32   `json:"id"`            // id
	Txt_Name      string  `json:"txt_Name"`      // 名称
	Type          int32   `json:"type"`          // 物品类型
	Color         int32   `json:"color"`         // 品质颜色
	Stack         int32   `json:"stack"`         // 堆叠上限
	UseLimit      int32   `json:"useLimit"`      // 使用上限
	Func          string  `json:"func"`          // 功能ID
	FuncParameter []int32 `json:"funcParameter"` // 功能参数
	BaseAttr      []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"baseAttr"` // 装备基础属性
	ArmorCompose []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"armorCompose"` // 装备合成与消耗
}

type itemTable struct {
	items map[int32]*item
}

func (self *itemTable) Load() {
	var arr []*item
	if !load_json("item.json", &arr) {
		return
	}

	items := make(map[int32]*item)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *itemTable) Query(id int32) *item {
	return self.items[id]
}

func (self *itemTable) Items() map[int32]*item {
	return self.items
}
