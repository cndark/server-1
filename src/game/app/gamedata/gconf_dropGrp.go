package gamedata

var ConfDropGrp = &dropGrpTable{}

type dropGrp struct {
	Seq     int32 `json:"seq"`     // 序列编号
	GroupId int32 `json:"groupId"` // 物品组ID
	Id      int32 `json:"id"`      // 物品ID
	Num     int64 `json:"num"`     // 物品数量
	Weight  int32 `json:"weight"`  // 掉落权重
}

type dropGrpTable struct {
	items map[int32]*dropGrp
}

func (self *dropGrpTable) Load() {
	var arr []*dropGrp
	if !load_json("dropGrp.json", &arr) {
		return
	}

	items := make(map[int32]*dropGrp)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *dropGrpTable) Query(seq int32) *dropGrp {
	return self.items[seq]
}

func (self *dropGrpTable) Items() map[int32]*dropGrp {
	return self.items
}
