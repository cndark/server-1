package guild

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/comp"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
)

// ============================================================================

type tech_t struct {
	Techs map[int32]int32 // [id]lv

	mod_job map[int32]comp.PropMods // [hero.job]prop

	plr IPlayer
}

// ============================================================================
func init() {
	// 给英雄挂上本模块的属性加成数据
	evtmgr.On(gconst.Evt_ModsHero, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		hero := args[1].(*comp.Hero)

		conf := gamedata.ConfMonster.Query(hero.Id)
		if conf != nil {

			m := plr.GetGuildPlrData().Tech.mod_job[conf.JobId]
			if m != nil {
				hero.GetMods().Add(m)
			}
		}
	})
}

func new_tech() *tech_t {
	return &tech_t{
		Techs: make(map[int32]int32),
	}
}

// ============================================================================

func (self *tech_t) init(plr IPlayer) {
	self.plr = plr
	self.mod_job = make(map[int32]comp.PropMods)

	self.init_mods()
}

func (self *tech_t) SetLevel(id int32, lv int32) {
	// conf
	conf := gamedata.ConfGuildTech.Query(id)
	if conf == nil {
		return
	}

	// find
	old := self.Techs[id]
	if lv == old {
		return
	}

	// set level
	self.Techs[id] = lv

	// calc props
	job, m := self.apply_mods(id, old, lv)
	self.apply_mods_heroes(job, m)
}

func (self *tech_t) Reset() {
	for _, h := range self.plr.GetBag().Heroes {
		conf := gamedata.ConfMonster.Query(h.Id)
		if conf != nil {

			m := self.mod_job[conf.JobId]
			if m != nil {
				h.GetMods().Sub(m)
				h.CalcProps(true)
			}
		}
	}

	// clear
	self.Techs = make(map[int32]int32)
	self.mod_job = make(map[int32]comp.PropMods)
}

// ============================================================================
// 属性

func (self *tech_t) init_mods() {
	for id, lv := range self.Techs {
		job, m := self.apply_mods(id, 0, lv)
		self.add_job_mod(job, m)
	}
}

// 本模块加
func (self *tech_t) add_job_mod(job int32, m comp.PropMods) {
	if job == 0 {
		return
	}

	jmod := self.mod_job[job]
	if jmod == nil {
		jmod = comp.NewPropMods()
		self.mod_job[job] = jmod
	}

	jmod.Add(m)
}

// 计数属性增量
func (self *tech_t) apply_mods(id int32, old_lv, new_lv int32) (job int32, ret comp.PropMods) {
	if old_lv == new_lv {
		return
	}

	// conf
	conf_tech := gamedata.ConfGuildTech.Query(id)
	if conf_tech == nil {
		return
	}

	job = conf_tech.Job
	ret = comp.NewPropMods()
	// old dec
	if old_lv > 0 {
		for _, v := range conf_tech.InitProp {
			amt := v.Val + v.Grow*float32(old_lv-1)
			ret.ModExt(v.Id, -amt)
		}
	}

	// add new
	if new_lv > 0 {
		for _, v := range conf_tech.InitProp {
			amt := v.Val + v.Grow*float32(new_lv-1)
			ret.ModExt(v.Id, amt)
		}
	}

	return
}

// 添加到本模块,和英雄身上
func (self *tech_t) apply_mods_heroes(job int32, m comp.PropMods) {
	self.add_job_mod(job, m)

	for _, h := range self.plr.GetBag().Heroes {
		conf := gamedata.ConfMonster.Query(h.Id)
		if conf != nil && conf.JobId == job {
			h.GetMods().Add(m)
			h.CalcProps(true)
		}
	}
}
