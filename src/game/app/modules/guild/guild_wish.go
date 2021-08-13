package guild

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"sync/atomic"
	"time"
)

// ============================================================================

var (
	seq_wish int64
)

// ============================================================================

type wish_t struct {
	Seq     int64
	PlrId   string
	ItemNum int32 // 许愿物编号
	Helps   int32 // 助力次数
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GuildLeave, func(args ...interface{}) {
		gld := args[0].(*Guild)
		plr := args[1].(IPlayer)

		gld.WishRemove(plr.GetId())
	})
}

// ============================================================================

func (self *Guild) WishCount(plrid string) (n int32) {
	for _, v := range self.Wish {
		if v.PlrId == plrid {
			n++
		}
	}
	return
}

func (self *Guild) WishAdd(plrid string, num int32) int64 {
	w := &wish_t{
		Seq:     time.Now().Unix()*10000 + atomic.AddInt64(&seq_wish, 1)%10000,
		PlrId:   plrid,
		ItemNum: num,
		Helps:   0,
	}

	self.Wish[w.Seq] = w

	return w.Seq
}

func (self *Guild) WishRemove(plrid string) {
	var arr []int64
	for _, v := range self.Wish {
		if v.PlrId == plrid {
			arr = append(arr, v.Seq)
		}
	}

	for _, seq := range arr {
		delete(self.Wish, seq)
	}
}
