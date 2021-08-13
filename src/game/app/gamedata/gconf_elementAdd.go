package gamedata

var ConfElementAdd = &elementAddTable{}

type elementAdd struct {
	Id      int32   `json:"id"`      // 序列
	Element []int32 `json:"element"` // 元素阵营
	Num     int32   `json:"num"`     // 人数
	Prop    []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"prop"` // 加成属性
}

type elementAddTable struct {
	items map[int32]*elementAdd
}

func (self *elementAddTable) Load() {
	var arr []*elementAdd
	if !load_json("elementAdd.json", &arr) {
		return
	}

	items := make(map[int32]*elementAdd)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *elementAddTable) Query(id int32) *elementAdd {
	return self.items[id]
}

func (self *elementAddTable) Items() map[int32]*elementAdd {
	return self.items
}
