package gamedata

var ConfVip = &vipTable{}

type vip struct {
	Lv      int32 `json:"lv"`  // VIP等级
	Exp     int32 `json:"exp"` // VIP升级经验
	Counter []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"counter"` // 特权次数，叠加
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 等级奖励
	MonthCardReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"monthCardReward"` // 月卡奖励
}

type vipTable struct {
	items map[int32]*vip
}

func (self *vipTable) Load() {
	var arr []*vip
	if !load_json("vip.json", &arr) {
		return
	}

	items := make(map[int32]*vip)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *vipTable) Query(lv int32) *vip {
	return self.items[lv]
}

func (self *vipTable) Items() map[int32]*vip {
	return self.items
}
