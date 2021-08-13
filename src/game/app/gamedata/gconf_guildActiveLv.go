package gamedata

var ConfGuildActiveLv = &guildActiveLvTable{}

type guildActiveLv struct {
	Lv           int32 `json:"lv"`           // 活跃等级
	NeedActScore int64 `json:"needActScore"` // 升级经验
	RewardDonate int32 `json:"rewardDonate"` // 个人获得贡献
	RewardFund   int32 `json:"rewardFund"`   // 家族获得功勋
}

type guildActiveLvTable struct {
	items map[int32]*guildActiveLv
}

func (self *guildActiveLvTable) Load() {
	var arr []*guildActiveLv
	if !load_json("guildActiveLv.json", &arr) {
		return
	}

	items := make(map[int32]*guildActiveLv)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *guildActiveLvTable) Query(lv int32) *guildActiveLv {
	return self.items[lv]
}

func (self *guildActiveLvTable) Items() map[int32]*guildActiveLv {
	return self.items
}
