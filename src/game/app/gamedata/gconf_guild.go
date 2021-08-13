package gamedata

var ConfGuild = &guildTable{}

type guild struct {
	GuildLevel    int32 `json:"guildLevel"` // 公会等级
	UpgradeExp    int64 `json:"upgradeExp"` // 升级经验
	MbLimit       int32 `json:"mbLimit"`    // 人数上限
	ViceLimit     int32 `json:"viceLimit"`  // 副会长人数
	CheckinReward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"checkinReward"` // 签到奖励
	CheckinExp int32 `json:"checkinExp"` // 签到经验
}

type guildTable struct {
	items map[int32]*guild
}

func (self *guildTable) Load() {
	var arr []*guild
	if !load_json("guild.json", &arr) {
		return
	}

	items := make(map[int32]*guild)

	for _, v := range arr {
		items[v.GuildLevel] = v
	}

	self.items = items
}

func (self *guildTable) Query(guildLevel int32) *guild {
	return self.items[guildLevel]
}

func (self *guildTable) Items() map[int32]*guild {
	return self.items
}
