package gamedata

var ConfActBillLtDayM = &actBillLtDayTableM{}

type actBillLtDayTableM struct {
	objs map[int32][]*actBillLtDay
}

func (self *actBillLtDayTableM) Load() {
	self.objs = make(map[int32][]*actBillLtDay)

	for _, v := range ConfActBillLtDay.Items() {
		key := v.ConfGrp
		self.objs[key] = append(self.objs[key], v)
	}
}

func (self *actBillLtDayTableM) QueryItems(key int32) []*actBillLtDay {
	return self.objs[key]
}
