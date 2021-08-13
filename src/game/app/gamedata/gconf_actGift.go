package gamedata

var ConfActGift = &actGiftTable{}

type actGift struct {
	Id      int32  `json:"id"`      // id
	ActName string `json:"actName"` // 活动
	Type    int32  `json:"type"`    // 礼包类型
	Reward  []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	GuildSharedReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"guildSharedReward"` // 公会共享奖励
	BuyCntLimit int32 `json:"buyCntLimit"` // 活动期间限购次数
	PayId       int32 `json:"payId"`       // payId
}

type actGiftTable struct {
	items map[int32]*actGift
}

func (self *actGiftTable) Load() {
	var arr []*actGift
	if !load_json("actGift.json", &arr) {
		return
	}

	items := make(map[int32]*actGift)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *actGiftTable) Query(id int32) *actGift {
	return self.items[id]
}

func (self *actGiftTable) Items() map[int32]*actGift {
	return self.items
}
