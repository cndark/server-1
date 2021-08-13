package gamedata

var ConfInviteReward = &inviteRewardTable{}

type inviteReward struct {
	Id        int32 `json:"id"` // 类型
	AttainTab []*struct {
		AttainId int32 `json:"attainId"`
		P2       int32 `json:"p2"`
	} `json:"attainTab"` // attain条件
	Reward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
}

type inviteRewardTable struct {
	items map[int32]*inviteReward
}

func (self *inviteRewardTable) Load() {
	var arr []*inviteReward
	if !load_json("inviteReward.json", &arr) {
		return
	}

	items := make(map[int32]*inviteReward)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *inviteRewardTable) Query(id int32) *inviteReward {
	return self.items[id]
}

func (self *inviteRewardTable) Items() map[int32]*inviteReward {
	return self.items
}
