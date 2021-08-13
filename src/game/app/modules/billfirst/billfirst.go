package billfirst

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

type BillFirst struct {
	Gear map[int32]*item

	plr IPlayer
}

type item struct {
	Day   int32   // 充值后登录天数
	Taken []int32 // 已领取的奖励
}

// ============================================================================

func NewBillFirst() *BillFirst {
	return &BillFirst{
		Gear: make(map[int32]*item),
	}
}

func init() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		prod_id := args[1].(int32)

		conf := gamedata.ConfBillProduct.Query(prod_id)
		if conf == nil {
			return
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_BillFirst)

		conf_data := gamedata.ConfBillFirst.Query(conf.PayId)
		if conf_data == nil {
			return
		}

		it := plr.GetBillFirst().Gear[conf_data.PayId]

		if it == nil {
			it = &item{Day: 1, Taken: []int32{}}
			plr.GetBillFirst().Gear[conf_data.PayId] = it
		} else {
			return
		}

		for _, v := range conf_data.Reward {
			op.Inc(v.Id, v.N)
		}

		op.Apply()

		plr.SendMsg(&msg.GS_BillFirstNew{
			Data: &msg.BillFirstItem{
				Id:    conf.PayId,
				Day:   it.Day,
				Taken: it.Taken,
			},
		})

	})

	evtmgr.On(gconst.Evt_PlrDailyOnline, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		for _, value := range plr.GetBillFirst().Gear {
			value.Day++
		}
	})
}

// ============================================================================

func (self *BillFirst) Init(plr IPlayer) {
	self.plr = plr
}

func (self *BillFirst) Take(id int32, day int32) (int32, *msg.Rewards) {

	titem := self.Gear[id]
	if titem == nil {
		return Err.BillFirst_NotBill, nil
	}

	conf := gamedata.ConfBillFirst.Query(id)
	if conf == nil {
		return Err.BillFirst_NotFound, nil
	}

	for _, v := range titem.Taken {
		if v == day {
			return Err.BillFirst_AlreadRewarded, nil
		}
	}

	if day > titem.Day {
		return Err.BillFirst_NotCond, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_BillFirst)
	for _, v := range conf.SignReward {
		if day == v.Day {
			op.Inc(v.Id, v.N)
		}
	}

	titem.Taken = append(titem.Taken, day)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

func (self *BillFirst) ToMsg() *msg.BillFirstData {

	ret := &msg.BillFirstData{}

	for id, v := range self.Gear {
		ret.Items = append(ret.Items, &msg.BillFirstItem{
			Id:    id,
			Day:   v.Day,
			Taken: v.Taken,
		})
	}

	return ret
}
