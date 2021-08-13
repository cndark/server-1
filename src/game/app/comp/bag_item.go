package comp

import (
	"fw/src/core/log"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

// ============================================================================

func (self *Bag) add_item(id int32, n int32) {
	if n == 0 {
		return
	}

	// check
	if !gconst.IsItem(id) {
		log.Warning("the id value is NOT item id:", id)
		return
	}

	conf := gamedata.ConfItem.Query(id)
	if conf == nil {
		log.Warning("item conf NOT found:", id)
		return
	}

	// add
	v := self.Items[id]
	v += n
	if conf.Stack > 0 && v > conf.Stack {
		v = conf.Stack
	}
	if v < 0 {
		v = 0
		log.Warning("item final count < 0:", id, n, v)
	}

	if v == 0 {
		delete(self.Items, id)
	} else {
		self.Items[id] = v
	}
}

func (self *Bag) GetItem(id int32) int32 {
	return self.Items[id]
}

func (self *Bag) ToMsg_ItemArray() (ret []*msg.Item) {
	for id, num := range self.Items {
		ret = append(ret, &msg.Item{
			Id:  id,
			Num: num,
		})
	}

	return
}
