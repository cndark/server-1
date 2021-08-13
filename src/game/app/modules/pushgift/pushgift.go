package pushgift

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	"strings"
	"time"
)

// ============================================================================
// 推送礼包
type PushGift struct {
	Gifts gift_m

	plr IPlayer
}

type gift_m map[int32]*gift_t
type gift_t struct {
	BuyCnt   int32
	CreateTs time.Time
}

type val_t float64

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		prod_id := args[1].(int32)
		csext := args[2].(string)

		prod := gamedata.ConfBillProduct.Query(prod_id)
		if prod == nil || prod.TypeId != gconst.Bill_Gift {
			return
		}

		arr := strings.Split(csext, "_")
		if len(arr) < 2 || arr[0] != gconst.Bill_CsExt_Type_PushGift {
			return
		}

		id := core.Atoi32(arr[1])
		plr.GetPushGift().bill_buy(id, prod.PayId)
	})
}

func NewPushGift() *PushGift {
	return &PushGift{
		Gifts: make(gift_m),
	}
}

// ============================================================================

func (self *PushGift) Init(plr IPlayer) {
	self.plr = plr
}

func (self *PushGift) isfin(id int32) bool {
	gt := self.Gifts[id]
	if gt != nil {
		return true
	}

	return false
}

func (self *PushGift) bill_buy(id int32, payid int32) {
	conf := gamedata.ConfPushGift.Query(id)
	if conf == nil || conf.PayId != payid {
		return
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_PushGift)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.plr.SendMsg(&msg.GS_PushGiftRewards{
		Id:      id,
		Rewards: rwds,
	})

	gift := self.Gifts[id]
	if gift == nil {
		gift = &gift_t{
			CreateTs: time.Now(),
		}
		self.Gifts[id] = gift
	}

	gift.BuyCnt++

	evtmgr.Fire(gconst.Evt_PushGift, self.plr, id)

}

func (self *PushGift) Get(id int32) *gift_t {
	gift := self.Gifts[id]
	if gift == nil {
		gift = &gift_t{}
		self.Gifts[id] = gift
	}

	return gift
}

func (self *PushGift) ToMsg() *msg.PushGiftData {
	ret := &msg.PushGiftData{}

	for id, v := range self.Gifts {
		ret.Gifts = append(ret.Gifts, &msg.PushGiftOne{
			Id:       id,
			BuyCnt:   v.BuyCnt,
			CreateTs: v.CreateTs.Unix(),
		})
	}

	return ret
}

// ============================================================================
// implements ICondVObj interface

func (self *val_t) GetVal() float64 {
	return float64(*self)
}

func (self *val_t) SetVal(v float64) {
	*self = val_t(v)
}

func (self *val_t) AddVal(v float64) {
	*self += val_t(v)
}

func (self *val_t) Done(body interface{}, confid int32, isChange bool) {
	conf := gamedata.ConfPushGift.Query(confid)
	if conf == nil {
		return
	}

	pg := body.(IPlayer).GetPushGift()
	if float64(*self) >= conf.P2 {
		gift := pg.Gifts[confid]
		if gift != nil {
			return
		}

		one := &gift_t{
			CreateTs: time.Now(),
		}

		pg.Gifts[confid] = one

		body.(IPlayer).SendMsg(&msg.GS_PushGiftNew{
			Id:       confid,
			CreateTs: one.CreateTs.Unix(),
		})
	}
}
