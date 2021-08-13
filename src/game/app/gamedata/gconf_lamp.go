package gamedata

var ConfLamp = &lampTable{}

type lamp struct {
	Id    int32   `json:"id"`    // 条件ID
	Type  int32   `json:"type"`  // 类型
	P1    int32   `json:"p1"`    // 参数1
	P2    []int32 `json:"p2"`    // 参数2
	Style int32   `json:"style"` // 跑马灯造型
}

type lampTable struct {
	items map[int32]*lamp
}

func (self *lampTable) Load() {
	var arr []*lamp
	if !load_json("lamp.json", &arr) {
		return
	}

	items := make(map[int32]*lamp)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *lampTable) Query(id int32) *lamp {
	return self.items[id]
}

func (self *lampTable) Items() map[int32]*lamp {
	return self.items
}
