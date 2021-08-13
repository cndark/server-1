package gamedata

var ConfGuildDonate = &guildDonateTable{}

type guildDonate struct {
	Id   int32 `json:"id"` // id
	Cost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"cost"` // 捐献消耗
	RewardExp    int32 `json:"rewardExp"`    // 家族获得经验
	RewardFund   int32 `json:"rewardFund"`   // 家族获得功勋
	RewardDonate int32 `json:"rewardDonate"` // 个人获得贡献
}

type guildDonateTable struct {
	items map[int32]*guildDonate
}

func (self *guildDonateTable) Load() {
	var arr []*guildDonate
	if !load_json("guildDonate.json", &arr) {
		return
	}

	items := make(map[int32]*guildDonate)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *guildDonateTable) Query(id int32) *guildDonate {
	return self.items[id]
}

func (self *guildDonateTable) Items() map[int32]*guildDonate {
	return self.items
}
