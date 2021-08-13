package gamedata

var ConfGuildBoss = &guildBossTable{}

type guildBoss struct {
	Num     int32 `json:"num"` // 层数
	Monster []*struct {
		Id int32 `json:"id"`
		Lv int32 `json:"lv"`
	} `json:"monster"` // boss
	RoundType       int32 `json:"roundType"` // 回合类型
	ChallengeReward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"challengeReward"` // 挑战奖励
	KillReward []*struct {
		A  int32 `json:"a"`
		B  int32 `json:"b"`
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"killReward"` // 击杀奖励
}

type guildBossTable struct {
	items map[int32]*guildBoss
}

func (self *guildBossTable) Load() {
	var arr []*guildBoss
	if !load_json("guildBoss.json", &arr) {
		return
	}

	items := make(map[int32]*guildBoss)

	for _, v := range arr {
		items[v.Num] = v
	}

	self.items = items
}

func (self *guildBossTable) Query(num int32) *guildBoss {
	return self.items[num]
}

func (self *guildBossTable) Items() map[int32]*guildBoss {
	return self.items
}
