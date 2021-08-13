package comp

import (
	"fw/src/core/log"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

// ============================================================================

func (self *Bag) add_ccy(id int32, n int64) {
	if n == 0 {
		return
	}

	// check
	if !gconst.IsCurrency(id) {
		log.Warning("the id value is NOT currency id:", id)
		return
	}

	// player exp interception
	if id == gconst.PlayerExp {
		self.plr.AddExp(int32(n))
		return
	}

	// vip exp interception
	if id == gconst.VipExp {
		self.plr.AddVipExp(int32(n))
		return
	}

	// add
	v := self.Ccy[id]
	v += n
	if v < 0 {
		v = 0
		log.Warning("currency final value < 0:", id, n, v)
	}

	if v == 0 {
		delete(self.Ccy, id)
	} else {
		self.Ccy[id] = v
	}
}

func (self *Bag) GetCcy(id int32) int64 {
	return self.Ccy[id]
}

func (self *Bag) ToMsg_CcyArray() (ret []*msg.Ccy) {
	for id, val := range self.Ccy {
		ret = append(ret, &msg.Ccy{
			Id:  id,
			Val: val,
		})
	}

	return
}
