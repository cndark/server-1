package mopen

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

// ============================================================================

// open模块
type MOpen struct {
	Modules m_mod_t `bson:"m"`

	plr IPlayer
}

type m_mod_t map[int32]*mod_t
type mod_t struct {
	Val  float64
	Open bool
}

// ============================================================================

func NewMOpen() *MOpen {
	return &MOpen{
		Modules: make(m_mod_t),
	}
}

func (self *MOpen) Init(plr IPlayer) {
	self.plr = plr

	conf := gamedata.ConfInitial.Query(1)
	for _, v := range conf.InitialModule {
		self.Modules[v] = &mod_t{Open: true}
	}
}

func (self *MOpen) get(id int32) (m *mod_t) {
	m = self.Modules[id]
	if m == nil {
		m = &mod_t{}
		self.Modules[id] = m
	}

	return
}

func (self *MOpen) Add(mid int32) {
	m := self.Modules[mid]
	if m == nil || m.Open {
		return
	}

	conf := gamedata.ConfOpen.Query(mid)
	if conf == nil {
		return
	}

	self.Modules[mid].Open = true

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_MOpen)
	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	// notify
	self.plr.SendMsg(&msg.GS_MOpenModuleNew{
		MId:     mid,
		Rewards: rwds,
	})

	// fire
	evtmgr.Fire(gconst.Evt_MOpen, self.plr, mid)
}

func (self *MOpen) IsOpen(mid int32) bool {
	m := self.Modules[mid]
	if m != nil {
		return m.Open
	}

	return false
}

func (self *MOpen) ToMsg() *msg.MOpenData {
	ret := &msg.MOpenData{
		M: make(map[int32]bool),
	}

	for id, v := range self.Modules {
		if v.Open {
			ret.M[id] = true
		}
	}

	return ret
}

// ============================================================================
// implements ICondVObj interface

func (self *mod_t) GetVal() float64 {
	return self.Val
}

func (self *mod_t) SetVal(v float64) {
	self.Val = v
}

func (self *mod_t) AddVal(v float64) {
	self.Val += v
}

func (self *mod_t) Done(body interface{}, confid int32, isChange bool) {
	conf := gamedata.ConfOpen.Query(confid)
	if conf == nil {
		return
	}

	mopen := body.(IPlayer).GetMOpen()

	if self.Val >= conf.P2 {
		mopen.Add(conf.ModuleId)
	}
}
