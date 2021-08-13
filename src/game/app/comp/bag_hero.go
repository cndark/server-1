package comp

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/worlddata"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
)

// ============================================================================

type Hero struct {
	Seq  int64 // 序号
	Id   int32 // 表 Id
	Lv   int32 // 等级
	Star int32 // 星级
	Skin int32 // 当前皮肤

	// 属性
	prop_map  PropMap
	prop_mods PropMods

	// 可装备栏位
	equip_box *equipbox_t

	// trinket
	Trinket *trinket_t

	// 锁定
	Locked bool

	// 即将转生id
	ChangeId int32

	// 战力
	atk_power int32

	bag *Bag
}

type trinket_t struct {
	Lv        int32
	Props     []int32
	props_tmp []int32
}

// ============================================================================

func new_hero(id int32) *Hero {
	// check
	if !gconst.IsHero(id) {
		log.Warning("the id value is NOT hero id:", id)
		return nil
	}

	conf := gamedata.ConfMonster.Query(id)
	if conf == nil {
		log.Warning("hero conf NOT found:", id)
		return nil
	}

	// new
	hero := &Hero{
		Seq:     worlddata.GenSeqHero(),
		Id:      id,
		Lv:      1,
		Star:    conf.Star,
		Trinket: &trinket_t{},
	}

	return hero
}

// ============================================================================

func (self *Bag) add_hero(hero *Hero) {
	if self.Heroes[hero.Seq] != nil {
		log.Warning("hero already added:", hero.Seq, hero.Id)
		return
	}

	self.Heroes[hero.Seq] = hero

	hero.init(self)
	hero.apply_mods(true)

	// fire
	if hero.Star >= 4 {
		evtmgr.Fire(gconst.Evt_HeroStar, self.plr, hero.Id, hero.Star)
	}
}

func (self *Bag) del_hero(seq int64) {
	hero := self.Heroes[seq]
	if hero == nil {
		return
	}

	// unequip all
	hero.NewEquipOp().UnequipAll().Apply()

	// delete
	delete(self.Heroes, seq)

	// cleanup
	hero.equip_box = nil
	hero.bag = nil
}

func (self *Bag) FindHero(seq int64) *Hero {
	return self.Heroes[seq]
}

func (self *Bag) ToMsg_HeroArray() (ret []*msg.Hero) {
	for _, hero := range self.Heroes {
		ret = append(ret, hero.ToMsg())
	}

	return
}

// 拥有指定等级的英雄个数
func (self *Bag) HeroReachLvCnt(lv int32) (n int32) {
	for _, hero := range self.Heroes {
		if hero.Lv >= lv {
			n++
		}
	}

	return
}

// 拥有指定星级的英雄个数
func (self *Bag) HeroReachStarCnt(star int32) (n int32) {
	for _, hero := range self.Heroes {
		if hero.Star >= star {
			n++
		}
	}

	return
}

// 最强英雄
func (self *Bag) MaxPowerHero() *Hero {
	var max int32
	var hero *Hero
	for _, v := range self.Heroes {
		if v.atk_power > int32(max) {
			max = v.atk_power
			hero = v
		}
	}

	return hero
}

// ============================================================================

func (self *Hero) init(bag *Bag) {
	self.prop_map = NewPropMap()
	self.prop_mods = NewPropMods()

	self.equip_box = new_equipbox(self)

	self.bag = bag
}

func (self *Hero) GetSeq() int64 {
	return self.Seq
}

func (self *Hero) GetElem() int32 {
	conf := gamedata.ConfMonster.Query(self.Id)
	if conf == nil {
		return 0
	} else {
		return conf.Elem
	}
}

func (self *Hero) apply_mods(send_update bool) {
	// 挂载全模块属性
	evtmgr.Fire(gconst.Evt_ModsHero, self.bag.plr, self)

	// 初始属性
	self.apply_mods_init()

	// 等级
	self.apply_mods_lv(0, self.Lv)

	// 星级
	self.apply_mods_star(0, self.Star)

	// 可装备栏
	self.apply_mods_equip()

	// trinket
	self.apply_mods_trinket(1)

	// calc
	self.CalcProps(send_update)
}

func (self *Hero) CalcProps(send_update bool) {
	// update prop-map
	updated := self.prop_map.Update(self.prop_mods)
	if len(updated) == 0 {
		return
	}

	// update atk-power
	self.update_atk_power()

	if send_update {
		// push: hero update
		self.send_update(updated)

		// fire
		evtmgr.Fire(gconst.Evt_HeroAtkPower, self.bag.plr, self.Seq, self.atk_power)
	}
}

func (self *Hero) GetMods() PropMods {
	return self.prop_mods
}

func (self *Hero) NewEquipOp() *EquipOp {
	return self.equip_box.new_op()
}

func (self *Hero) SetLevel(lv int32) {
	if lv == self.Lv {
		return
	}

	// set new level
	old := self.Lv
	self.Lv = lv

	// re-calc-props
	self.apply_mods_lv(old, lv)
	self.CalcProps(true)

	// fire
	evtmgr.Fire(gconst.Evt_HeroLv, self.bag.plr, self.Id, lv)
}

func (self *Hero) SetStar(star int32) {
	if star == self.Star {
		return
	}

	// set new Star
	old := self.Star
	self.Star = star

	// re-calc-props
	self.apply_mods_star(old, star)
	self.CalcProps(true)

	// fire
	evtmgr.Fire(gconst.Evt_HeroStar, self.bag.plr, self.Id, star)
}

func (self *Hero) SetChangeId(id int32) {
	self.ChangeId = id
}

func (self *Hero) ChangeApply() {
	if self.ChangeId == 0 {
		return
	}

	self.Id = self.ChangeId
	self.ChangeId = 0

	self.prop_map = NewPropMap()
	self.prop_mods = NewPropMods()
	self.apply_mods(true)
}

func (self *Hero) Inherit(id int32) {
	if id == self.Id {
		return
	}

	self.Id = id

	self.prop_map = NewPropMap()
	self.prop_mods = NewPropMods()
	self.apply_mods(true)
}

func (self *Hero) SetSkin(id int32) {
	self.Skin = id
}

func (self *Hero) update_atk_power() {
	atk_power := float32(0)
	rat := float32(0)

	// from props
	for _, attr := range gamedata.ConfAttribute.Items() {
		if attr.PowerType == 1 {
			atk_power += self.prop_map.Get(attr.AttributeId) * attr.HeroPower
		} else if attr.PowerType == 2 {
			rat += self.prop_map.Get(attr.AttributeId) * attr.HeroPower
		}
	}

	atk_power *= (1 + rat)

	// modify power
	conf_m := gamedata.ConfMonster.Query(self.Id)
	if conf_m != nil {
		atk_power -= float32(conf_m.ModifyPower)
	}

	// set
	self.atk_power = int32(atk_power)
}

func (self *Hero) GetAtkPower() int32 {
	return self.atk_power
}

func (self *Hero) send_update(updated PropMap) {
	hero := self.ToMsg()
	hero.Props = map[int32]float32(updated)

	self.bag.plr.SendMsg(&msg.GS_HeroUpdate{
		Hero: hero,
	})
}

func (self *Hero) apply_mods_init() {
	conf := gamedata.ConfMonster.Query(self.Id)
	if conf == nil {
		return
	}

	for _, v := range conf.HeroBaseProps {
		self.prop_mods.ModBase(v.Id, v.Val)
	}
}

func (self *Hero) apply_mods_lv(old_lv, new_lv int32) {
	conf := gamedata.ConfMonster.Query(self.Id)
	if conf == nil {
		return
	}

	// remove old mods
	if old_lv > 1 {
		for _, v := range conf.HeroPropGrowth {
			self.prop_mods.ModBase(v.Id, -v.Val*float32(old_lv-1))
		}
	}

	// add new mods
	if new_lv > 1 {
		for _, v := range conf.HeroPropGrowth {
			self.prop_mods.ModBase(v.Id, v.Val*float32(new_lv-1))
		}
	}
}

func (self *Hero) apply_mods_star(old_star, new_star int32) {
	// remove old mods
	{
		conf := gamedata.ConfHeroStarUp.Query(old_star)
		if conf != nil {
			for _, v := range conf.PropsRatio {
				self.prop_mods.ModExt(v.Id, -v.Val)
			}
		}
	}

	// add new mods
	{
		conf := gamedata.ConfHeroStarUp.Query(new_star)
		if conf != nil {
			for _, v := range conf.PropsRatio {
				self.prop_mods.ModExt(v.Id, v.Val)
			}
		}
	}
}

func (self *Hero) apply_mods_equip() {
	mods := self.equip_box.get_mods()
	self.prop_mods.Add(mods)
}

// ============================================================================
// trinket

func (self *Hero) TrinketUnlock() int32 {
	if self.Lv < 40 {
		return Err.Hero_LowLevel
	}
	if self.Trinket.Lv > 0 {
		return Err.Hero_TrinketAlreadyUnlocked
	}

	self.Trinket.Lv = 1
	self.trinket_gen_rand_props(false)
	self.trinket_commit_rand_props()
	self.apply_mods_trinket(1)
	self.CalcProps(true)

	return Err.OK
}

func (self *Hero) TrinketSetLevel(lv int32, lock bool) {
	if self.Trinket.Lv == lv {
		return
	}

	self.apply_mods_trinket(-1)

	self.Trinket.Lv = lv
	self.trinket_gen_rand_props(lock)
	self.trinket_commit_rand_props()

	self.apply_mods_trinket(1)
	self.CalcProps(true)

	evtmgr.Fire(gconst.Evt_HeroTrinketLv, self.bag.plr, self.Seq, lv)
}

func (self *Hero) TrinketTransformGen() []int32 {
	self.trinket_gen_rand_props(false)
	return self.Trinket.props_tmp
}

func (self *Hero) TrinketTransformCommit() int32 {
	if self.Trinket.props_tmp == nil {
		return Err.Hero_TrinketNoPropToCommit
	}

	self.apply_mods_trinket(-1)
	self.trinket_commit_rand_props()
	self.apply_mods_trinket(1)
	self.CalcProps(true)

	return Err.OK
}

func (self *Hero) trinket_commit_rand_props() {
	self.Trinket.Props = self.Trinket.props_tmp
	self.Trinket.props_tmp = nil
}

func (self *Hero) trinket_gen_rand_props(lock bool) {
	// conf
	conf := gamedata.ConfTrinket.Query(self.Trinket.Lv)
	if conf == nil {
		return
	}

	// make initial prop-lib
	roll := make([]*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
		Wt  int32   `json:"wt"`
	}, 0, len(conf.Prop))

	// empty tmp props
	self.Trinket.props_tmp = make([]int32, 0, conf.AttributeNum)

	// check lock
	N := len(self.Trinket.Props)
	left := conf.AttributeNum

	if lock {
		for _, id := range self.Trinket.Props {
			self.Trinket.props_tmp = append(self.Trinket.props_tmp, id)
		}
		left -= int32(N)
	}

	sum := int32(0)
	for i, v := range conf.Prop {
		if lock && core.ArrayFind(N, func(j int) bool { return self.Trinket.Props[j] == v.Id }) >= 0 {
			continue
		}

		roll = append(roll, conf.Prop[i])
		sum += v.Wt
	}

	// roll
	for i := int32(0); i < left && sum > 0; i++ {
		// random pick
		p := rand.Int31n(sum)
		for i, v := range roll {
			p -= v.Wt
			if p < 0 {
				// pick
				self.Trinket.props_tmp = append(self.Trinket.props_tmp, v.Id)

				// remove entry
				L := len(roll)
				roll[i] = roll[L-1]
				roll = roll[:L-1]

				// update wt sum
				sum -= v.Wt

				break
			}
		}
	}
}

func (self *Hero) apply_mods_trinket(sign int) {
	conf := gamedata.ConfTrinket.Query(self.Trinket.Lv)
	if conf == nil {
		return
	}

	for _, id := range self.Trinket.Props {
		i := core.ArrayFind(len(conf.Prop), func(i int) bool {
			return conf.Prop[i].Id == id
		})
		if i >= 0 {
			self.prop_mods.ModExt(id, conf.Prop[i].Val*float32(sign))
		}
	}
}

// ============================================================================

func (self *Hero) ToMsg() *msg.Hero {
	return &msg.Hero{
		Seq:  self.Seq,
		Id:   self.Id,
		Lv:   self.Lv,
		Star: self.Star,
		Trinket: &msg.Trinket{
			Lv:    self.Trinket.Lv,
			Props: self.Trinket.Props,
		},
		Props:    self.prop_map,
		Locked:   self.Locked,
		AtkPwr:   self.atk_power,
		ChangeId: self.ChangeId,
		Skin:     self.Skin,
	}
}

func (self *Hero) ToMsg_Detail() *msg.HeroDetail {
	ret := &msg.HeroDetail{
		Hero:  self.ToMsg(),
		Relic: make(map[int32]int32),
	}

	for _, v := range self.equip_box.box {
		if v == nil {
			continue
		}

		id := v.equip_get_id()
		if gconst.IsItem(id) {
			conf := gamedata.ConfItem.Query(id)
			if conf != nil && gconst.IsArmor(conf.Type) {
				ret.Armors = append(ret.Armors, id)
			}
		} else if gconst.IsRelic(id) {
			relic := self.bag.FindRelic(v.equip_get_seq())
			if relic != nil {
				ret.Relic[id] = relic.Star
			}
		}
	}

	return ret
}

func (self *Hero) ToMsg_BattleFighter(pos int32) *msg.BattleFighter {
	return &msg.BattleFighter{
		Seq:    self.Seq,
		Id:     self.Id,
		Lv:     self.Lv,
		Star:   self.Star,
		Props:  self.prop_map,
		Pos:    pos,
		AtkPwr: self.atk_power,
		Skin:   self.Skin,
	}
}
