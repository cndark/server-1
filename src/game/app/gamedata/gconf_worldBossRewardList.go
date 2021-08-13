package gamedata

var ConfWorldBossRewardList = &worldBossRewardListTable{}

type worldBossRewardList struct {
	Id     int32 `json:"id"` // 奖励id
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 阶段奖励
	BuffId []int32 `json:"buffId"` // 加成buffId
}

type worldBossRewardListTable struct {
	items map[int32]*worldBossRewardList
}

func (self *worldBossRewardListTable) Load() {
	var arr []*worldBossRewardList
	if !load_json("worldBossRewardList.json", &arr) {
		return
	}

	items := make(map[int32]*worldBossRewardList)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *worldBossRewardListTable) Query(id int32) *worldBossRewardList {
	return self.items[id]
}

func (self *worldBossRewardListTable) Items() map[int32]*worldBossRewardList {
	return self.items
}
