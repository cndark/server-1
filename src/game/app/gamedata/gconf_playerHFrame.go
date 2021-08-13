package gamedata

var ConfPlayerHFrame = &playerHFrameTable{}

type playerHFrame struct {
	Id     int32 `json:"id"` // id
	Status []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"status"` // 前置条件
	Cost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 激活消耗道具
	Time int32 `json:"time"` // 时限(天)
}

type playerHFrameTable struct {
	items map[int32]*playerHFrame
}

func (self *playerHFrameTable) Load() {
	var arr []*playerHFrame
	if !load_json("playerHFrame.json", &arr) {
		return
	}

	items := make(map[int32]*playerHFrame)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *playerHFrameTable) Query(id int32) *playerHFrame {
	return self.items[id]
}

func (self *playerHFrameTable) Items() map[int32]*playerHFrame {
	return self.items
}
