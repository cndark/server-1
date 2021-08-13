package gamedata

var ConfLanguage = &languageTable{}

type language struct {
	Key  string `json:"key"`  // 语言id
	Text string `json:"text"` // 名称
}

type languageTable struct {
	items map[string]*language
}

func (self *languageTable) Load() {
	var arr []*language
	if !load_json("language.json", &arr) {
		return
	}

	items := make(map[string]*language)

	for _, v := range arr {
		items[v.Key] = v
	}

	self.items = items
}

func (self *languageTable) Query(key string) *language {
	return self.items[key]
}

func (self *languageTable) Items() map[string]*language {
	return self.items
}
