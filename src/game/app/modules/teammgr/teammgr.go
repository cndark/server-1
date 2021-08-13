package teammgr

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

// ============================================================================

// 阵容管理
type TeamMgr struct {
	Teams map[int32]*msg.TeamFormation

	plr IPlayer
}

// ============================================================================

func init() {
	// update player atk-power by hero atkpower
	evtmgr.On(gconst.Evt_HeroAtkPower, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		seq := args[1].(int64)

		if !plr.GetTeamMgr().InTeam(seq, gconst.TeamType_Dfd) {
			return
		}

		evtmgr.Fire(gconst.Evt_DfdTeamUpdate, plr, plr.GetTeamMgr().GetTeam(gconst.TeamType_Dfd))
	})

	evtmgr.On(gconst.Evt_HeroDel, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		seqs := args[1].([]int64)

		for tp, t := range plr.GetTeamMgr().Teams {
			if tp == gconst.TeamType_Dfd {
				continue
			}

			if t.Formation != nil {
				for _, seq := range seqs {
					_, ok := t.Formation[seq]
					if ok {
						delete(t.Formation, seq)
					}
				}
			}
		}
	})
}

func NewTeamMgr() *TeamMgr {
	return &TeamMgr{
		Teams: make(map[int32]*msg.TeamFormation),
	}
}

// ============================================================================

func (self *TeamMgr) Init(plr IPlayer) {
	self.plr = plr
}

func (self *TeamMgr) GetTeam(tp int32) *msg.TeamFormation {
	t := self.Teams[tp]
	if t == nil || len(t.Formation) == 0 {
		t = self.Teams[gconst.TeamType_Dfd]
	}

	return t
}

func (self *TeamMgr) SetTeam(tp int32, T *msg.TeamFormation) {
	if tp <= 0 || tp >= gconst.TeamType_Max {
		return
	}

	self.Teams[tp] = T

	if tp == gconst.TeamType_Dfd {
		evtmgr.Fire(gconst.Evt_DfdTeamUpdate, self.plr, self.GetTeam(tp))
	}
}

func (self *TeamMgr) IsSetTeam(tp int32) bool {
	t := self.Teams[tp]
	if t == nil {
		return false
	}

	return len(t.Formation) > 0
}

func (self *TeamMgr) InTeam(seq int64, tps ...int32) bool {
	if len(tps) == 0 {
		for _, t := range self.Teams {
			_, ok := t.Formation[seq]
			if ok {
				return true
			}
		}
	} else {
		for _, tp := range tps {
			t := self.Teams[tp]
			if t != nil {
				_, ok := t.Formation[seq]
				if ok {
					return true
				}
			}
		}
	}

	return false
}

// 获取阵容战力
func (self *TeamMgr) TeamAtkPwr(tp int32) (v int32) {
	t := self.Teams[tp]
	if t != nil {
		for seq := range t.Formation {
			hero := self.plr.GetBag().FindHero(seq)
			if hero == nil {
				continue
			}

			v += hero.GetAtkPower()
		}
	}

	return
}

// 只求普通防守阵容战力
func (self *TeamMgr) CalcPlrAtkPwr() int32 {
	return self.TeamAtkPwr(gconst.TeamType_Dfd)
}

// ============================================================================

func (self *TeamMgr) ToMsg() *msg.TeamMgrData {
	return &msg.TeamMgrData{
		Teams: self.Teams,
	}
}
