package gamedata

var ConfAppointM = &appointTableM{}

type appointTableM struct {
	objs map[int32][]*appoint
}

func (self *appointTableM) Load() {
	self.objs = make(map[int32][]*appoint)

	for _, v := range ConfAppoint.Items() {
		key := v.Star
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *appointTableM) QueryItems(star int32) []*appoint {
	return self.objs[star]
}
