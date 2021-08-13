package gamedata

var ConfMailM = &mailTableM{}

type mailTableM struct {
	items map[string]*mail
}

func (self *mailTableM) Load() {
	self.items = make(map[string]*mail)

	for _, v := range ConfMail.items {
		if v.Skey != "" {
			self.items[v.Skey] = v
		}
	}
}

func (self *mailTableM) Query(skey string) *mail {
	return self.items[skey]
}
