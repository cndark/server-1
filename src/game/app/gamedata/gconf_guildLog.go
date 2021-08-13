package gamedata

var ConfGuildLog = &guildLogTable{}

type guildLog struct {
	Id int32 `json:"id"` // id
}

type guildLogTable struct {
	items map[int32]*guildLog
}

func (self *guildLogTable) Load() {
	var arr []*guildLog
	if !load_json("guildLog.json", &arr) {
		return
	}

	items := make(map[int32]*guildLog)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *guildLogTable) Query(id int32) *guildLog {
	return self.items[id]
}

func (self *guildLogTable) Items() map[int32]*guildLog {
	return self.items
}
