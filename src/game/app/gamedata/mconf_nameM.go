package gamedata

var ConfNameM = &nameTableM{}

type nameTableM struct {
	objs map[string][]*name
}

func (self *nameTableM) Load() {
	self.objs = make(map[string][]*name)

	for _, v := range ConfName.Items() {
		key := v.Language
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *nameTableM) QueryItems(key string) []*name {
	return self.objs[key]
}
