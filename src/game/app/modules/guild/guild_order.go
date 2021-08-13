package guild

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"sync/atomic"
	"time"
)

// ============================================================================

var (
	seq_order int64
)

// ============================================================================

type order_t struct {
	Orders     map[int64]*order_rec_t
	GetOrderTs time.Time

	plr IPlayer
}

type order_rec_t struct {
	Seq     int64
	Star    int32
	StartTs time.Time
}

// ============================================================================

func new_order() *order_t {
	return &order_t{
		Orders: make(map[int64]*order_rec_t),
	}
}

func (self *order_t) init(plr IPlayer) {
	self.plr = plr
}

func (self *order_t) GetOrders() (ec int32, ret []*msg.GuildOrderRec) {
	// get gld
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	// conf
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return Err.Failed, nil
	}

	// cd
	now := time.Now()
	if now.Sub(self.GetOrderTs).Hours() < float64(conf.HarborOrderCd) {
		return Err.Guild_OrderCd, nil
	}

	// how many orders can we get
	conf_hb := gamedata.ConfGuildHarbor.Query(gld.Harbor.Lv)
	if conf_hb == nil {
		return Err.Failed, nil
	}

	n := int(conf_hb.OrderLimit)
	left := int(conf.GuildOrderLimit) - len(self.Orders)
	if n > left {
		n = left
	}
	if n <= 0 {
		return Err.Guild_OrderLimit, nil
	}

	// get n orders
	ret = make([]*msg.GuildOrderRec, 0, n)
	for i := 0; i < n; i++ {
		o := self.roll_order(gld.Harbor.Lv)
		if o == nil {
			continue
		}

		self.Orders[o.Seq] = o
		ret = append(ret, &msg.GuildOrderRec{
			Seq:     o.Seq,
			Star:    o.Star,
			StartTs: o.StartTs.Unix(),
		})
	}

	self.GetOrderTs = now
	return Err.OK, ret
}

func (self *order_t) roll_order(lv int32) *order_rec_t {
	// sum
	sum := int32(0)
	for _, v := range gamedata.ConfGuildOrder.Items() {
		sum += v.Weight
	}
	if sum <= 0 {
		return nil
	}

	// roll
	p := rand.Int31n(sum)
	for _, v := range gamedata.ConfGuildOrder.Items() {
		p -= v.Weight
		if p < 0 {
			return &order_rec_t{
				Seq:  time.Now().Unix()*10000 + atomic.AddInt64(&seq_order, 1)%10000,
				Star: 1,
			}
		}
	}

	return nil // should NOT happen
}

func (self *order_t) Starup(seq int64) int32 {
	// get gld
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember
	}

	// find order
	o := self.Orders[seq]
	if o == nil {
		return Err.Guild_OrderNotFound
	}

	// is started ?
	if !o.StartTs.IsZero() {
		return Err.Guild_OrderAlreadyStarted
	}

	// full star ?
	conf := gamedata.ConfGuildOrder.Query(o.Star)
	if conf == nil || len(conf.UpStarCost) == 0 {
		return Err.Guild_OrderFullStar
	}

	// cost
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GuildOrderStarup)
	for _, v := range conf.UpStarCost {
		op.Dec(v.Id, v.N)
	}
	if ec := op.CheckEnough(); ec != Err.OK {
		return ec
	}

	op.Apply()

	// up
	o.Star++

	return Err.OK
}

func (self *order_t) StartOrder(seq int64) (ec int32, ts int64) {
	// get gld
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, 0
	}

	// find order
	o := self.Orders[seq]
	if o == nil {
		return Err.Guild_OrderNotFound, 0
	}

	// is started ?
	if !o.StartTs.IsZero() {
		return Err.Guild_OrderAlreadyStarted, 0
	}

	// ok. start it
	o.StartTs = time.Now()

	return Err.OK, o.StartTs.Unix()
}

func (self *order_t) CloseOrder(seq int64) (ec int32, rwd *msg.Rewards) {
	// get gld
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	// find order
	o := self.Orders[seq]
	if o == nil {
		return Err.Guild_OrderNotFound, nil
	}

	// conf
	conf_o := gamedata.ConfGuildOrder.Query(o.Star)
	if conf_o == nil {
		return Err.Failed, nil
	}

	conf_h := gamedata.ConfGuildHarbor.Query(gld.Harbor.Lv)
	if conf_h == nil {
		return Err.Failed, nil
	}

	// ended ?
	need_dur := float64(conf_o.Time) * float64(1-conf_h.TimeReduce)
	if o.StartTs.IsZero() || time.Since(o.StartTs).Minutes() < need_dur {
		return Err.Guild_OrderNotEnd, nil
	}

	// close it
	delete(self.Orders, seq)

	// rewards
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GuildOrderClose)
	for _, v := range conf_o.Reward {
		op.Inc(v.Id, float32(v.N)*(1+conf_h.OutAdd))
	}
	rwd = op.Apply().ToMsg()

	evtmgr.Fire(gconst.Evt_GuildOrderClose, self.plr, conf_o.Star)

	return Err.OK, rwd
}

func (self *order_t) OrderList() (ret []*msg.GuildOrderRec) {
	// get gld
	gld := self.plr.GetGuild()
	if gld == nil {
		return
	}

	// list
	ret = make([]*msg.GuildOrderRec, 0, len(self.Orders))
	for _, v := range self.Orders {
		ret = append(ret, &msg.GuildOrderRec{
			Seq:     v.Seq,
			Star:    v.Star,
			StartTs: v.StartTs.Unix(),
		})
	}

	return
}
