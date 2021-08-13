package gamedata

var ConfActBillLtTotal = &actBillLtTotalTable{}

type actBillLtTotal struct {
	Seq      int32 `json:"seq"`      // 序列
	BillCond int32 `json:"billCond"` // 充值基准货币值条件
	Reward   []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	ConfGrp int32 `json:"confGrp"` // 活动配置
}

type actBillLtTotalTable struct {
	items map[int32]*actBillLtTotal
}

func (self *actBillLtTotalTable) Load() {
	var arr []*actBillLtTotal
	if !load_json("actBillLtTotal.json", &arr) {
		return
	}

	items := make(map[int32]*actBillLtTotal)

	for _, v := range arr {
		items[v.Seq] = v
	}

	self.items = items
}

func (self *actBillLtTotalTable) Query(seq int32) *actBillLtTotal {
	return self.items[seq]
}

func (self *actBillLtTotalTable) Items() map[int32]*actBillLtTotal {
	return self.items
}
