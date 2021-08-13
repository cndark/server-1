package gamedata

var ConfGuildShrineLv = &guildShrineLvTable{}

type guildShrineLv struct {
	Lv     int32 `json:"lv"` // 神灵等级
	Status []*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	} `json:"status"` // 家族等级要求
	NeedShrineScore int32 `json:"needShrineScore"` // 升级所需经验
	RewardFund      int32 `json:"rewardFund"`      // 获得家族功勋（财富）
	RewardElem      int32 `json:"rewardElem"`      // 获得魂玉（奢侈品）
	RewardDonate    int32 `json:"rewardDonate"`    // 获得个人贡献
}

type guildShrineLvTable struct {
	items map[int32]*guildShrineLv
}

func (self *guildShrineLvTable) Load() {
	var arr []*guildShrineLv
	if !load_json("guildShrineLv.json", &arr) {
		return
	}

	items := make(map[int32]*guildShrineLv)

	for _, v := range arr {
		items[v.Lv] = v
	}

	self.items = items
}

func (self *guildShrineLvTable) Query(lv int32) *guildShrineLv {
	return self.items[lv]
}

func (self *guildShrineLvTable) Items() map[int32]*guildShrineLv {
	return self.items
}
