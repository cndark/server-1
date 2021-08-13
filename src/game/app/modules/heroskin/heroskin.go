package heroskin

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/comp"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 英雄皮肤库
type HeroSkin struct {
	Skins map[int32]*hero_skin_t

	mod_hero map[int32]comp.PropMods // 同id英雄才加

	plr IPlayer
}

type hero_skin_t struct {
	Lv int32 // 等级
}

// ============================================================================

func init() {
	// 给英雄挂上本模块的属性加成数据
	evtmgr.On(gconst.Evt_ModsHero, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		hero := args[1].(*comp.Hero)

		m := plr.GetHeroSkin().mod_hero[hero.Id]
		hero.GetMods().Add(m)
	})
}

func NewHeroSkin() *HeroSkin {
	return &HeroSkin{
		Skins: make(map[int32]*hero_skin_t),
	}
}

// ============================================================================

func (self *HeroSkin) Init(plr IPlayer) {
	self.plr = plr
	self.mod_hero = make(map[int32]comp.PropMods)

	self.init_mods()
}

func (self *HeroSkin) IsExist(id int32) bool {
	return self.Skins[id] != nil
}

// 激活皮肤
func (self *HeroSkin) AddSkin(id int32) int32 {
	conf := gamedata.ConfHeroSkin.Query(id)
	if conf == nil {
		return Err.Failed
	}

	if self.IsExist(id) {
		return Err.Hero_SkinExist
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_HeroSkin)
	for _, v := range conf.Cost {
		op.Dec(v.Id, v.N)
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec
	}

	self.Skins[id] = &hero_skin_t{
		Lv: 1,
	}

	op.Apply()

	// mods
	m := comp.NewPropMods()
	for _, v := range conf.BaseProps {
		m.ModExt(v.Id, v.Val)
	}
	self.apply_mods_heroes(conf.Hero, m)

	return Err.OK
}

// 皮肤升级
func (self *HeroSkin) LevelUp(id int32) int32 {
	conf := gamedata.ConfHeroSkin.Query(id)
	if conf == nil {
		return Err.Failed
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return Err.Failed
	}

	if !self.IsExist(id) {
		return Err.Hero_SkinNotExist
	}

	skin := self.Skins[id]
	if skin.Lv >= conf_g.SkinLvLimit {
		return Err.Hero_SkinLvFull
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_HeroSkin)
	for _, v := range conf.LvUpCost {
		op.Dec(v.Id, v.N)
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec
	}

	skin.Lv++

	op.Apply()

	// mods
	m := comp.NewPropMods()
	for _, v := range conf.BasePropGrowth {
		m.ModExt(v.Id, v.Val)
	}
	self.apply_mods_heroes(conf.Hero, m)

	return Err.OK
}

// ============================================================================

func (self *HeroSkin) ToMsg() *msg.HeroSkinData {
	ret := &msg.HeroSkinData{
		Skins: make(map[int32]int32),
	}

	for id, v := range self.Skins {
		ret.Skins[id] = v.Lv
	}

	return ret
}

// ============================================================================
// 属性

// 初始属性
func (self *HeroSkin) init_mods() {
	for id, skin := range self.Skins {
		m := self.apply_mods(id, skin.Lv)

		conf := gamedata.ConfHeroSkin.Query(id)
		if conf != nil {
			self.add_hero_mod(conf.Hero, m)
		}
	}
}

// 本模块加
func (self *HeroSkin) add_hero_mod(heroid int32, m comp.PropMods) {
	hmod := self.mod_hero[heroid]
	if hmod == nil {
		hmod = comp.NewPropMods()
		self.mod_hero[heroid] = hmod
	}

	hmod.Add(m)
}

// 计数属性
func (self *HeroSkin) apply_mods(id int32, lvN int32) (ret comp.PropMods) {
	conf := gamedata.ConfHeroSkin.Query(id)
	if conf == nil {
		return
	}

	skin := self.Skins[id]
	if skin == nil {
		return
	}

	ret = comp.NewPropMods()
	for _, v := range conf.BaseProps {
		ret.ModExt(v.Id, v.Val)
	}

	for _, v := range conf.BasePropGrowth {
		ret.ModExt(v.Id, v.Val*float32(lvN))
	}

	return
}

// 添加到本模块,和英雄身上
func (self *HeroSkin) apply_mods_heroes(heroid int32, m comp.PropMods) {
	self.add_hero_mod(heroid, m)

	for _, h := range self.plr.GetBag().Heroes {
		if h.Id == heroid {
			h.GetMods().Add(m)
			h.CalcProps(true)
		}
	}
}
