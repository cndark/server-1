package gamedata

var ConfEmotion = &emotionTable{}

type emotion struct {
	Id int32 `json:"id"` // 表情ID
}

type emotionTable struct {
	items map[int32]*emotion
}

func (self *emotionTable) Load() {
	var arr []*emotion
	if !load_json("emotion.json", &arr) {
		return
	}

	items := make(map[int32]*emotion)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *emotionTable) Query(id int32) *emotion {
	return self.items[id]
}

func (self *emotionTable) Items() map[int32]*emotion {
	return self.items
}
