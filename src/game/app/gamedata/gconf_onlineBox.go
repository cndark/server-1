package gamedata

var ConfOnlineBox = &onlineBoxTable{}

type onlineBox struct {
	Id         int32 `json:"id"`         // 序列
	OnlineTime int32 `json:"onlineTime"` // 在线时长要求(分钟)
	Reward     []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
}

type onlineBoxTable struct {
	items map[int32]*onlineBox
}

func (self *onlineBoxTable) Load() {
	var arr []*onlineBox
	if !load_json("onlineBox.json", &arr) {
		return
	}

	items := make(map[int32]*onlineBox)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *onlineBoxTable) Query(id int32) *onlineBox {
	return self.items[id]
}

func (self *onlineBoxTable) Items() map[int32]*onlineBox {
	return self.items
}
