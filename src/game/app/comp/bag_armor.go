package comp

import (
	"fw/src/core/log"
	. "fw/src/core/math"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/worlddata"
	"fw/src/game/msg"
	"sort"
)

// ============================================================================

type Armor struct {
	Seq int64 // 序号
	Id  int32 // 表 Id

	HeroSeq int64 // 穿戴英雄Seq

	bag *Bag
}

// ============================================================================

func (self *Bag) PrepareArmor(id int32) *Armor {
	// check
	if !gconst.IsItem(id) {
		log.Warning("the id value is NOT item id:", id)
		return nil
	}

	conf := gamedata.ConfItem.Query(id)
	if conf == nil || !gconst.IsArmor(conf.Type) {
		log.Warning("NOT an armor id:", id)
		return nil
	}

	return &Armor{
		Id:  id,
		bag: self,
	}
}

func (self *Bag) FindArmor(seq int64) *Armor {
	return self.Armors[seq]
}

func (self *Bag) ToMsg_ArmorArray() (ret []*msg.Armor) {
	for _, ar := range self.Armors {
		ret = append(ret, ar.ToMsg())
	}

	return
}

// ============================================================================

func (self *Armor) init(bag *Bag) {
	self.bag = bag

	// add to hero equip-box
	if self.HeroSeq > 0 {
		hero := self.bag.FindHero(self.HeroSeq)
		if hero != nil {
			hero.equip_box.add(self)
		}
	}
}

func (self *Armor) apply_mods_prop(mods PropMods) {
	conf := gamedata.ConfItem.Query(self.Id)
	if conf == nil {
		return
	}

	for _, v := range conf.BaseAttr {
		mods.ModExt(v.Id, v.Val)
	}
}

func (self *Armor) ToMsg() *msg.Armor {
	return &msg.Armor{
		Seq:     self.Seq,
		Id:      self.Id,
		HeroSeq: self.HeroSeq,
	}
}

// ============================================================================
// 可穿戴注册

func init() {
	equip_register_group(&equip_group_t{
		id:              gconst.EquipGroup_Armor,
		slot_start:      gconst.EquipGroup_Armor_SlotStart,
		slot_end:        gconst.EquipGroup_Armor_SlotEnd,
		apply_mods_func: armor_apply_mods_group,
	})
}

func armor_apply_mods_group(box []IEquippable, a, b int32, mods PropMods) {
	// single armor mods
	for i := a; i <= b; i++ {
		e := box[i]
		if e != nil {
			e.(*Armor).apply_mods_prop(mods)
		}
	}

	// armor master mods
	{
		lvs := make([]int32, 0, 4)
		for i := a; i <= b; i++ {
			e := box[i]
			if e != nil {
				lvs = append(lvs, e.(*Armor).Id%100)
			}
		}

		sort.Slice(lvs, func(i, j int) bool {
			return lvs[i] > lvs[j]
		})

		N := MinInt(4, len(lvs))
		for i := 2; i <= N; i++ {
			conf := gamedata.ConfArmorMaster.Query(lvs[i-1])
			if conf != nil {
				p := conf.Prop[i-2]
				mods.ModExt(p.Id, p.Val)
			}
		}
	}
}

// ============================================================================
// 可穿戴接口

func (self *Armor) equip_get_grpid() int {
	return gconst.EquipGroup_Armor
}

func (self *Armor) equip_get_id() int32 {
	return self.Id
}

func (self *Armor) equip_get_seq() int64 {
	return self.Seq
}

func (self *Armor) eqiup_is_valid() bool {
	return self.bag.Items[self.Id] > 0
}

func (self *Armor) equip_get_slot(box []IEquippable) int32 {
	conf := gamedata.ConfItem.Query(self.Id)
	if conf == nil {
		return 0
	} else {
		return gconst.ArmorSlot(conf.Type)
	}
}

func (self *Armor) equip_get_bind_seq() int64 {
	return self.HeroSeq
}

func (self *Armor) equip_set_bind_seq(seq int64, slot int32) {
	self.HeroSeq = seq

	if seq > 0 {
		self.Seq = worlddata.GenSeqArmor()

		// remove from item container
		self.bag.Items[self.Id]--

		if self.bag.Items[self.Id] <= 0 {
			delete(self.bag.Items, self.Id)
		}

		// add to armor container
		self.bag.Armors[self.Seq] = self

	} else if seq == 0 {
		// remove from armor container
		delete(self.bag.Armors, self.Seq)

		// add to item container
		self.bag.Items[self.Id]++
	}

	// notify
	self.bag.plr.SendMsg(&msg.GS_ArmorUpdate_HeroSeq{
		Seq:     self.Seq,
		Id:      self.Id,
		HeroSeq: self.HeroSeq,
	})
}
