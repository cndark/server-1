package gamedata

var ConfPushGift = &pushGiftTable{}

type pushGift struct {
	Id     int32   `json:"id"`   // id
	Cond   int32   `json:"cond"` // cond
	P1     int32   `json:"p1"`   // p1
	P2     float64 `json:"p2"`   // p2
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	PayId       int32 `json:"payId"`       // 充值项
	BuyCntLimit int32 `json:"buyCntLimit"` // （推送时）限购次数
}

type pushGiftTable struct {
	items map[int32]*pushGift
}

func (self *pushGiftTable) Load() {
	var arr []*pushGift
	if !load_json("pushGift.json", &arr) {
		return
	}

	items := make(map[int32]*pushGift)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *pushGiftTable) Query(id int32) *pushGift {
	return self.items[id]
}

func (self *pushGiftTable) Items() map[int32]*pushGift {
	return self.items
}
