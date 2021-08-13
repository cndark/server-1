package comp

import (
	"fw/src/core/log"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/worlddata"
	"fw/src/game/msg"
)

// ============================================================================

type Relic struct {
	Seq  int64 // 序号
	Id   int32 // 表 Id
	Star int32
	Xp   int32

	HeroSeq int64 // 穿戴英雄Seq

	bag *Bag
}

// ============================================================================

func new_relic(id int32) *Relic {
	// check
	if !gconst.IsRelic(id) {
		log.Warning("the id value is NOT relic id:", id)
		return nil
	}

	conf := gamedata.ConfRelic.Query(id)
	if conf == nil {
		log.Warning("relic conf NOT found:", id)
		return nil
	}

	// new
	rlc := &Relic{
		Seq:  worlddata.GenSeqRelic(),
		Id:   id,
		Star: 1,
		Xp:   0,
	}

	return rlc
}

func (self *Bag) add_relic(rlc *Relic) {
	if self.Relics[rlc.Seq] != nil {
		log.Warning("relic already added:", rlc.Seq, rlc.Id)
		return
	}

	self.Relics[rlc.Seq] = rlc

	rlc.init(self)
}

func (self *Bag) del_relic(seq int64) {
	rlc := self.Relics[seq]
	if rlc == nil {
		return
	}

	delete(self.Relics, seq)

	// unequip if it's on a hero
	if rlc.HeroSeq > 0 {
		hero := self.FindHero(rlc.HeroSeq)
		if hero != nil {
			hero.NewEquipOp().Unequip(gconst.EquipGroup_Relic_SlotStart).Apply()
		}
	}

	rlc.bag = nil
}

func (self *Bag) FindRelic(seq int64) *Relic {
	return self.Relics[seq]
}

func (self *Bag) ToMsg_RelicArray() (ret []*msg.Relic) {
	for _, rlc := range self.Relics {
		ret = append(ret, rlc.ToMsg())
	}

	return
}

// ============================================================================

func (self *Relic) init(bag *Bag) {
	self.bag = bag

	// add to hero equip-box
	if self.HeroSeq > 0 {
		hero := self.bag.FindHero(self.HeroSeq)
		if hero != nil {
			hero.equip_box.add(self)
		}
	}
}

func (self *Relic) apply_mods_prop(mods PropMods) {
	// conf
	conf := gamedata.ConfRelic.Query(self.Id)
	if conf == nil {
		return
	}

	// base
	for _, v := range conf.BaseProps {
		mods.ModExt(v.Id, v.Val)
	}

	// growth
	for _, v := range conf.BasePropGrowth {
		mods.ModExt(v.Id, v.Val*(float32(self.Star-1)))
	}

	// hidden
	if self.HeroSeq > 0 {
		// check activation
		b := false
		for {
			hero := self.bag.FindHero(self.HeroSeq)
			if hero == nil {
				break
			}

			conf_hero := gamedata.ConfMonster.Query(hero.Id)
			if conf_hero == nil {
				break
			}

			if len(conf.Active) == 0 {
				break
			}

			if conf.Active[0].Type == 1 && conf.Active[0].N == conf_hero.Elem ||
				conf.Active[0].Type == 2 && conf.Active[0].N == conf_hero.JobId ||
				conf.Active[0].Type == 3 && conf.Active[0].N == conf_hero.Id {
				b = true
				break
			}

			break
		}

		if b {
			// hidden base
			for _, v := range conf.HideProps {
				mods.ModExt(v.Id, v.Val)
			}

			// hidden growth
			for _, v := range conf.HidePropGrowth {
				mods.ModExt(v.Id, v.Val*(float32(self.Star-1)))
			}
		}
	}
}

func (self *Relic) AddExp(v int32) {
	if v <= 0 {
		return
	}

	conf := gamedata.ConfRelic.Query(self.Id)
	if conf == nil {
		return
	}

	star := self.Star
	self.Xp += v

	// consume xp
	for {
		// check if star full
		if star >= conf.StarLimit {
			self.Xp = 0
			break
		}

		// check remaining xp
		need_xp := conf.Exp * star
		if self.Xp < need_xp {
			break
		}

		self.Xp -= need_xp
		star++
	}

	// check starup
	if star == self.Star {
		// only xp changes
		self.bag.plr.SendMsg(&msg.GS_RelicUpdate_Star{
			Seq:  self.Seq,
			Star: -1, // 无变化
			Xp:   self.Xp,
		})
	} else {
		// star up
		self.bag.plr.SendMsg(&msg.GS_RelicUpdate_Star{
			Seq:  self.Seq,
			Star: star,
			Xp:   self.Xp,
		})

		// apply mods if equipped
		hero := self.bag.FindHero(self.HeroSeq)
		if hero != nil {
			hero.NewEquipOp().Modify(self, func() {
				self.Star = star
			}).Apply()
		}

		// update star
		self.Star = star
	}
}

func (self *Relic) ToMsg() *msg.Relic {
	return &msg.Relic{
		Seq:     self.Seq,
		Id:      self.Id,
		Star:    self.Star,
		Xp:      self.Xp,
		HeroSeq: self.HeroSeq,
	}
}

// ============================================================================
// 可穿戴注册

func init() {
	equip_register_group(&equip_group_t{
		id:              gconst.EquipGroup_Relic,
		slot_start:      gconst.EquipGroup_Relic_SlotStart,
		slot_end:        gconst.EquipGroup_Relic_SlotEnd,
		apply_mods_func: relic_apply_mods_group,
	})
}

func relic_apply_mods_group(box []IEquippable, a, b int32, mods PropMods) {
	for i := a; i <= b; i++ {
		e := box[i]
		if e != nil {
			e.(*Relic).apply_mods_prop(mods)
		}
	}
}

// ============================================================================
// 可穿戴接口

func (self *Relic) equip_get_grpid() int {
	return gconst.EquipGroup_Relic
}

func (self *Relic) equip_get_id() int32 {
	return self.Id
}

func (self *Relic) equip_get_seq() int64 {
	return self.Seq
}

func (self *Relic) eqiup_is_valid() bool {
	return true
}

func (self *Relic) equip_get_slot(box []IEquippable) int32 {
	return gconst.EquipGroup_Relic_SlotStart
}

func (self *Relic) equip_get_bind_seq() int64 {
	return self.HeroSeq
}

func (self *Relic) equip_set_bind_seq(seq int64, slot int32) {
	self.HeroSeq = seq

	// notify
	self.bag.plr.SendMsg(&msg.GS_RelicUpdate_HeroSeq{
		Seq:     self.Seq,
		HeroSeq: self.HeroSeq,
	})
}
