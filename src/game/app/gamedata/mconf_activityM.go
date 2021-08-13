package gamedata

import (
	"fw/src/core"
	"strings"
)

var ConfActivityM = &activityTableM{}

type activityM struct {
	*activity
	ConfCycles []*activity_cycle_t
}

type activity_cycle_t struct {
	Cnt  int32
	Grps []int32
}

type activityTableM struct {
	items map[string]*activityM
}

func (self *activityTableM) Load() {
	self.items = make(map[string]*activityM)

	for _, itm := range ConfActivity.items {
		cycles := make([]*activity_cycle_t, 0, len(itm.ConfCycles))

		for _, c := range itm.ConfCycles {
			arr := strings.Split(c.Grps, ",")
			keys := make([]int32, 0, len(arr))
			for _, v := range arr {
				keys = append(keys, core.Atoi32(v))
			}

			cycles = append(cycles, &activity_cycle_t{
				Cnt:  c.Cnt,
				Grps: keys,
			})
		}

		self.items[itm.Name] = &activityM{
			activity:   itm,
			ConfCycles: cycles,
		}
	}
}

func (self *activityTableM) Query(name string) *activityM {
	return self.items[name]
}

func (self *activityTableM) Items() map[string]*activityM {
	return self.items
}
