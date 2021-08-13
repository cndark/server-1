package gamedata

var ConfCondDrop = &condDropTable{}

type condDrop struct {
	Id       int32  `json:"id"`       // 掉落ID
	Txt_Desc string `json:"txt_Desc"` // 描述（程序不读）
	Status   []*struct {
		StatusId int32 `json:"statusId"`
		Val      int32 `json:"val"`
	} `json:"status"` // 开启条件
	Items []*struct {
		Id   int32 `json:"id"`
		N    int64 `json:"n"`
		Prob int32 `json:"prob"`
	} `json:"items"` // 增量掉落
}

type condDropTable struct {
	items map[int32]*condDrop
}

func (self *condDropTable) Load() {
	var arr []*condDrop
	if !load_json("condDrop.json", &arr) {
		return
	}

	items := make(map[int32]*condDrop)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *condDropTable) Query(id int32) *condDrop {
	return self.items[id]
}

func (self *condDropTable) Items() map[int32]*condDrop {
	return self.items
}
