package gamedata

var ConfPlayerUp = &playerUpTable{}

type playerUp struct {
	Level    int32 `json:"level"` // 等级
	Exp      int32 `json:"exp"`   // 升级经验
	UpReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"upReward"` // 升级奖励
	FriendLimit   int32 `json:"friendLimit"` // 好友数量限制
	RiftMineRatio []*struct {
		Id int32   `json:"id"`
		N  float32 `json:"n"`
	} `json:"riftMineRatio"` // 裂隙矿产奖励系数
	RiftMonsterRatio []*struct {
		Id int32   `json:"id"`
		N  float32 `json:"n"`
	} `json:"riftMonsterRatio"` // 裂隙怪物奖励系数
	RiftBoxRatio []*struct {
		Id int32   `json:"id"`
		N  float32 `json:"n"`
	} `json:"riftBoxRatio"` // 裂隙宝箱奖励系数
}

type playerUpTable struct {
	items map[int32]*playerUp
}

func (self *playerUpTable) Load() {
	var arr []*playerUp
	if !load_json("playerUp.json", &arr) {
		return
	}

	items := make(map[int32]*playerUp)

	for _, v := range arr {
		items[v.Level] = v
	}

	self.items = items
}

func (self *playerUpTable) Query(level int32) *playerUp {
	return self.items[level]
}

func (self *playerUpTable) Items() map[int32]*playerUp {
	return self.items
}
