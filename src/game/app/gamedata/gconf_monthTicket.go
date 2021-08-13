package gamedata

var ConfMonthTicket = &monthTicketTable{}

type monthTicket struct {
	Lv         int32 `json:"lv"`     // 等级
	Exp        int32 `json:"exp"`    // 升级需要经验
	UpCost     int32 `json:"upCost"` // 直升消耗钻石
	BaseReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"baseReward"` // 基础奖励
	TicketReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"ticketReward"` // 月票奖励
}

type monthTicketTable struct {
	items map[int32]*monthTicket
}

func (self *monthTicketTable) Load() {
	var arr []*monthTicket
	if !load_json("monthTicket.json", &arr) {
		return
	}

	items := make(map[int32]*monthTicket)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *monthTicketTable) Query(lv int32) *monthTicket {
	return self.items[lv]
}

func (self *monthTicketTable) Items() map[int32]*monthTicket {
	return self.items
}
