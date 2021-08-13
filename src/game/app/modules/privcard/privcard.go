package privcard

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 特权卡
type PrivCard struct {
	Cards priv_card_m

	plr IPlayer
}

type priv_card_m map[int32]*priv_card_t
type priv_card_t struct {
	Id       int32
	ExpireTs time.Time
	Counter  map[int32]int64
	IsAward  bool
	AddCnt   int32
}

// ============================================================================

func init() {
	// reset daily
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		plr.GetPrivCard().daily_reset()
	})

	// bill done
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		pid := args[1].(int32)

		conf_prod := gamedata.ConfBillProduct.Query(pid)
		if conf_prod == nil || conf_prod.TypeId != gconst.Bill_PrivCard {
			return
		}

		for _, items := range gamedata.ConfPrivCard.Items() {
			if items.DisPayId == conf_prod.PayId ||
				items.PayId == conf_prod.PayId {

				plr.GetPrivCard().AddPrivCard(items.Id)
				break
			}
		}
	})
}

func NewPrivCard() *PrivCard {
	return &PrivCard{
		Cards: make(priv_card_m),
	}
}

// ============================================================================

func (self *PrivCard) Init(plr IPlayer) {
	self.plr = plr

	self.init_counter()
}

func (self *PrivCard) init_counter() {
	for _, card := range self.Cards {
		for id, n := range card.Counter {
			self.plr.GetCounter().Add(id, n, true)
		}

		conf := gamedata.ConfPrivCard.Query(card.Id)
		if conf != nil {
			for _, v := range conf.CounterOnce {
				self.plr.GetCounter().Add(v.Id, int64(v.N), true)
			}
		}
	}
}

func (self *PrivCard) daily_reset() {
	now := time.Now()
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_PrivCard)

	for _, card := range self.Cards {
		if !card.ExpireTs.After(now) {
			for id, n := range card.Counter {
				cop.DecCounterMax(id, n)
			}

			// clear
			card.Counter = make(map[int32]int64)
		}

		card.IsAward = false
	}

	cop.Apply()
}

func (self *PrivCard) IsPrivCardValid(id int32) bool {
	pc := self.Cards[id]
	if pc == nil {
		return false
	}

	if pc.ExpireTs.After(time.Now()) {
		return true
	}

	return false
}

func (self *PrivCard) AddPrivCard(id int32) {
	conf := gamedata.ConfPrivCard.Query(id)
	if conf == nil {
		return
	}

	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_PrivCard)

	for _, v := range conf.Reward {
		cop.Inc(v.Id, v.N)
	}

	zero := core.StartOfDay(time.Now())
	addTs := time.Duration(conf.ExpireDay) * 24 * time.Hour

	c := self.Cards[id]
	if c == nil {
		c = &priv_card_t{
			Id:       id,
			ExpireTs: zero.Add(addTs),
			Counter:  make(map[int32]int64),
		}

		self.Cards[id] = c

		for _, v := range conf.CounterTimeLimit {
			cop.IncCounterMax(v.Id, int64(v.N))
			c.Counter[v.Id] += int64(v.N)
		}

		for _, v := range conf.CounterOnce {
			cop.IncCounterMax(v.Id, int64(v.N))
		}

	} else if !c.ExpireTs.After(zero) {
		for _, v := range conf.CounterTimeLimit {
			cop.IncCounterMax(v.Id, int64(v.N))
			c.Counter[v.Id] += int64(v.N)
		}

		if c.AddCnt >= conf.ExtNeedAddCnt { // 额外奖励
			addTs = time.Duration(conf.ExpireDay+conf.ExtDay) * 24 * time.Hour

			for _, v := range conf.ExtReward {
				cop.Inc(v.Id, v.N)
			}

			for _, v := range conf.ExtCounter {
				cop.IncCounterMax(v.Id, int64(v.N))
				c.Counter[v.Id] += int64(v.N)
			}
		}

		c.ExpireTs = zero.Add(addTs)

	} else { // 有效期内重复获得
		if c.AddCnt >= conf.ExtNeedAddCnt { // 额外奖励
			addTs = time.Duration(conf.ExpireDay+conf.ExtDay) * 24 * time.Hour

			for _, v := range conf.ExtReward {
				cop.Inc(v.Id, v.N)
			}

			if c.AddCnt == conf.ExtNeedAddCnt { // 首次进入额外获得
				for _, v := range conf.ExtCounter {
					cop.IncCounterMax(v.Id, int64(v.N))
					c.Counter[v.Id] += int64(v.N)
				}
			}
		}

		c.ExpireTs = c.ExpireTs.Add(addTs)
	}

	c.AddCnt++

	rwds := cop.Apply().ToMsg()

	// notify
	self.plr.SendMsg(&msg.GS_PrivCardNew{
		Card: &msg.PrivCard{
			Id:       c.Id,
			ExpireTs: c.ExpireTs.Unix(),
			IsAward:  c.IsAward,
			AddCnt:   c.AddCnt,
		},
		Rewards: rwds,
	})

	// fire
	evtmgr.Fire(gconst.Evt_PrivCardAdd, self.plr, id)

	return
}

// 领取每日奖励
func (self *PrivCard) DailyTake(id int32) (int32, *msg.Rewards) {
	pc := self.Cards[id]
	if pc == nil {
		return Err.Plr_PrivCardNotFound, nil
	}

	if pc.ExpireTs.Before(time.Now()) {
		return Err.Plr_PrivCardNotFound, nil
	}

	if pc.IsAward {
		return Err.Plr_TakenBefore, nil
	}

	conf := gamedata.ConfPrivCard.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_PrivCard)
	for _, v := range conf.DailyReward {
		op.Inc(v.Id, v.N)
	}

	if id == gconst.C_PrivCard_Month {
		conf_vip := gamedata.ConfVip.Query(self.plr.GetVipLevel())
		if conf_vip != nil {
			for _, v := range conf_vip.MonthCardReward {
				op.Inc(v.Id, v.N)
			}
		}
	}

	pc.IsAward = true

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// ============================================================================

func (self *PrivCard) ToMsg() *msg.PrivCardData {
	ret := &msg.PrivCardData{}

	for _, v := range self.Cards {
		ret.Cards = append(ret.Cards, &msg.PrivCard{
			Id:       v.Id,
			ExpireTs: v.ExpireTs.Unix(),
			IsAward:  v.IsAward,
			AddCnt:   v.AddCnt,
		})
	}

	return ret
}
