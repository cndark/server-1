package gamedata

var ConfActBillLtDay = &actBillLtDayTable{}

type actBillLtDay struct {
	Seq     int32 `json:"seq"`     // 序列
	BillDay int32 `json:"billDay"` // 充值天数
	Reward  []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	ConfGrp int32 `json:"confGrp"` // 活动配置
}

type actBillLtDayTable struct {
	items map[int32]*actBillLtDay
}

func (self *actBillLtDayTable) Load() {
	var arr []*actBillLtDay
	if !load_json("actBillLtDay.json", &arr) {
		return
	}

	items := make(map[int32]*actBillLtDay)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actBillLtDayTable) Query(seq int32) *actBillLtDay {
	return self.items[seq]
}

func (self *actBillLtDayTable) Items() map[int32]*actBillLtDay {
	return self.items
}
