package vip

import (
	"fw/src/core/evtmgr"
	"fw/src/core/sched/async"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

type Vip struct {
	Lv  int32
	Exp int32

	plr IPlayer
}

// ============================================================================
func init() {

	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		prod_id := args[1].(int32)

		prod := gamedata.ConfBillProduct.Query(prod_id)
		if prod == nil {
			return
		}

		plr.GetVip().AddExp(prod.VipExp)
	})
}

func NewVip() *Vip {
	return &Vip{
		Lv: 0,
	}
}

// ============================================================================

func (self *Vip) Init(plr IPlayer) {
	self.plr = plr

	self.init_counter()
}

func (self *Vip) init_counter() {
	for i := int32(1); i <= self.Lv; i++ {
		conf := gamedata.ConfVip.Query(i)
		if conf != nil {
			for _, v := range conf.Counter {
				self.plr.GetCounter().Add(v.Id, int64(v.N), true)
			}
		}
	}
}

func (self *Vip) AddExp(v int32) {
	if v <= 0 {
		return
	}

	new_lv := self.Lv
	self.Exp += v

	for {
		conf := gamedata.ConfVip.Query(new_lv)
		if conf == nil || conf.Exp == 0 {
			break
		}

		if self.Exp < conf.Exp {
			break
		}

		new_lv++
	}

	if self.Lv != new_lv {
		n := new_lv - self.Lv
		self.Lv = new_lv
		self.on_level_up(n)
	}

	self.plr.SendMsg(&msg.GS_VipUpdate{
		Lv:  self.Lv,
		Exp: self.Exp,
	})
}

func (self *Vip) LevelUp() int32 {
	conf := gamedata.ConfVip.Query(self.Lv)
	if conf == nil || conf.Exp == 0 {
		return Err.Failed
	}

	self.Lv++
	self.Exp = 0

	self.on_level_up(1)

	self.plr.SendMsg(&msg.GS_VipUpdate{
		Lv:  self.Lv,
		Exp: self.Exp,
	})

	return Err.OK
}

func (self *Vip) on_level_up(n int32) {
	if n < 0 {
		return
	}

	// update userinfo
	async.Push(
		func() {
			dbmgr.Share_UpdateUserVip(self.plr.GetId(), self.Lv)
		},
	)

	// rewards
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	conf_m := gamedata.ConfMail.Query(conf.VipUpRewardMailId)
	if conf_m == nil {
		return
	}

	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_Vip)
	m := mail.New(self.plr).SetKey(conf.VipUpRewardMailId)

	for i := int32(0); i < n; i++ {
		conf := gamedata.ConfVip.Query(self.Lv - i)
		if conf != nil {
			for _, v := range conf.Counter {
				cop.IncCounterMax(v.Id, int64(v.N))
			}

			for _, v := range conf.Reward {
				m.AddAttachment(v.Id, float64(v.N))
			}
		}
	}

	cop.Apply()
	m.Send()

	// fire
	evtmgr.Fire(gconst.Evt_VipLv, self.plr, self.Lv)
}

// ============================================================================

func (self *Vip) ToMsg() *msg.VipData {
	return &msg.VipData{
		Lv:  self.Lv,
		Exp: self.Exp,
	}
}
