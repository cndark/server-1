package gamedata

var ConfBillFirst = &billFirstTable{}

type billFirst struct {
	PayId  int32 `json:"payId"` // 支付id
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 充值奖励
	SignReward []*struct {
		Day int32 `json:"day"`
		Id  int32 `json:"id"`
		N   int64 `json:"n"`
	} `json:"signReward"` // 签到获得
}

type billFirstTable struct {
	items map[int32]*billFirst
}

func (self *billFirstTable) Load() {
	var arr []*billFirst
	if !load_json("billFirst.json", &arr) {
		return
	}

	items := make(map[int32]*billFirst)

	for _, v := range arr {
		items[v.PayId] = v
	}

	self.items = items
}

func (self *billFirstTable) Query(payId int32) *billFirst {
	return self.items[payId]
}

func (self *billFirstTable) Items() map[int32]*billFirst {
	return self.items
}
