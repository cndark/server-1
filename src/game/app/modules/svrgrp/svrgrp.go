package svrgrp

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"fw/src/shared/config"
)

// ============================================================================

var (
	s_grp_cross *Group // 区服组
	s_grp_all   *Group // 全服组
)

// ============================================================================

type Group struct {
	Id       int32   // grpid
	Svrs     []int32 // full members
	AvailLen int32   // current available length of 'svrs'
	LBLen    int32   // load balance length of 'svrs'. updated at night (5:00 am)
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_ConfReload, func(args ...interface{}) {
		s_grp_cross.update_avail_len()

		s_grp_all.Svrs = config.GameIds
		s_grp_all.update_avail_len()
	})
}

// ============================================================================

func Init() {
	make_grp_cross()
	make_grp_all()
	sched_lb_len_update()

	evtmgr.Fire(gconst.Evt_SvrGrpReady)
}

func make_grp_cross() {
	// find which group we're in.
	// and make the group
	for _, v := range gamedata.ConfCrossconf.Items() {
		if core.ArrayFind(
			len(v.Svrs),
			func(i int) bool { return v.Svrs[i] == config.CurGame.Id },
		) >= 0 {
			s_grp_cross = &Group{
				Id:   v.Id,
				Svrs: v.Svrs,
			}

			s_grp_cross.update_avail_len()
			s_grp_cross.update_lb_len()

			return
		}
	}

	// not found
	core.Panic("game is NOT in any of the server-group:", config.CurGame.Id)
}

func make_grp_all() {
	s_grp_all = &Group{
		Id:   0,
		Svrs: config.GameIds,
	}

	s_grp_all.update_avail_len()
	s_grp_all.update_lb_len()
}

func sched_lb_len_update() {
	// 5:00 am
	u, t := core.ParseRepeatTime("d/5")

	var f func()
	f = func() {
		t = core.AddTimeByUnit(t, u, 1)
		loop.SetTimeout(t, func() {
			s_grp_cross.update_lb_len()
			s_grp_all.update_lb_len()
			f()
		})
	}
	f()
}

// ============================================================================
// group

func GetGroupCross() *Group {
	return s_grp_cross
}

func GetGroupAll() *Group {
	return s_grp_all
}

func (self *Group) update_avail_len() {
	maxid := config.GameIdMax

	self.AvailLen = 0
	L := len(self.Svrs)
	for i := L - 1; i >= 0; i-- {
		if self.Svrs[i] <= maxid {
			self.AvailLen = int32(i) + 1
			break
		}
	}
}

func (self *Group) update_lb_len() {
	self.LBLen = self.AvailLen
}

func (self *Group) IsMaster() bool {
	return config.CurGame.Id == self.Svrs[0]
}

func (self *Group) Has(svrid int32) bool {
	return core.ArrayFind(
		int(self.AvailLen),
		func(i int) bool { return self.Svrs[i] == svrid },
	) >= 0
}

func (self *Group) Broadcast(message msg.Message) {
	for _, id := range self.Svrs[:self.AvailLen] {
		utils.Send2Game(id, message)
	}
}

func (self *Group) Send2Master(message msg.Message) {
	utils.Send2Game(self.Svrs[0], message)
}

func (self *Group) Send2Player(svrid int32, plrid string, message msg.Message) {
	if svrid == config.CurGame.Id {
		utils.Send2Player(plrid, message)
	} else {
		utils.Send2CrossPlayer(svrid, plrid, message)
	}
}

func (self *Group) Push2Master(evt_name string, sarg []string, oarg interface{}) {
	utils.GsPush(self.Svrs[0], evt_name, sarg, oarg)
}

func (self *Group) PushAll(evt_name string, sarg []string, oarg interface{}, except_self ...bool) {
	// marshal object arg
	var oarg_data []byte
	var err error
	if oarg != nil {
		oarg_data, err = utils.MarshalArg(oarg)
		if err != nil {
			return
		}
	}

	// push all
	b := core.DefFalse(except_self)
	curid := config.CurGame.Id

	for _, id := range self.Svrs[:self.AvailLen] {
		if b && id == curid {
			continue
		}

		utils.Send2Game(id, &msg.GS_Push{
			EvtName: evt_name,
			SArg:    sarg,
			OArg:    oarg_data,
		})
	}
}
