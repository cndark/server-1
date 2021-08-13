package comp

import (
	"fw/src/core/log"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

const (
	eop_Equip   = 1
	eop_Unequip = 2
	eop_Modify  = 3
)

// ============================================================================

var (
	equip_group_registry = make(map[int]*equip_group_t) // [grpid]
	equip_slot_max       int32
)

// ============================================================================

type equipbox_t struct {
	owner IEquipOwner
	box   []IEquippable
}

type IEquipOwner interface {
	GetSeq() int64
	GetMods() PropMods
	CalcProps(send_update bool)
}

type IEquippable interface {
	equip_get_grpid() int                     // return equippable group id
	equip_get_id() int32                      // return equippable conf id
	equip_get_seq() int64                     // return equippable seq
	eqiup_is_valid() bool                     // is valid when equipping
	equip_get_slot(box []IEquippable) int32   // the slot when equipping or current slot
	equip_get_bind_seq() int64                //
	equip_set_bind_seq(seq int64, slot int32) // should include id-update to client
}

type EquipOp struct {
	ops  []*equip_op_entry_t // ordered op list
	ebox *equipbox_t         // the box
	err  int32               // last error
}

type equip_op_entry_t struct {
	tp   int         // op type
	e    IEquippable // equip itself
	slot int32       // slot info
	f    func()      // modify func
}

type equip_group_t struct {
	id              int
	slot_start      int32
	slot_end        int32
	apply_mods_func func(box []IEquippable, a, b int32, ret PropMods)
}

// ============================================================================

func equip_register_group(grp *equip_group_t) {
	equip_group_registry[grp.id] = grp
	if grp.slot_end > equip_slot_max {
		equip_slot_max = grp.slot_end
	}
}

func equip_is_slot_valid(e IEquippable, slot int32) bool {
	if e == nil {
		return slot >= 0 && slot <= equip_slot_max
	} else {
		grp := equip_group_registry[e.equip_get_grpid()]
		return grp != nil && slot >= grp.slot_start && slot <= grp.slot_end
	}
}

// ============================================================================

func new_equipbox(owner IEquipOwner) *equipbox_t {
	return &equipbox_t{
		owner: owner,
		box:   make([]IEquippable, equip_slot_max+1),
	}
}

func (self *equipbox_t) get(slot int32) IEquippable {
	if equip_is_slot_valid(nil, slot) {
		return self.box[slot]
	} else {
		return nil
	}
}

func (self *equipbox_t) add(e IEquippable) {
	box := self.box

	slot := e.equip_get_slot(box)
	if !equip_is_slot_valid(nil, slot) {
		return
	}

	if box[slot] == nil {
		box[slot] = e
	} else {
		e.equip_set_bind_seq(0, slot)
		log.Warning("equip box already has an item in slot:", slot)
	}
}

func (self *equipbox_t) new_op() *EquipOp {
	return &EquipOp{
		ops:  nil,
		ebox: self,
		err:  Err.OK,
	}
}

func (self *equipbox_t) __equip(e IEquippable, slot int32) int32 {
	box := self.box

	// check validation
	if !e.eqiup_is_valid() {
		return Err.Equip_NotFound
	}

	// check if it's equipped
	if e.equip_get_bind_seq() > 0 {
		return Err.Equip_AlreadyEquipped
	}

	// do it
	if box[slot] == nil {
		// equip
		box[slot] = e
		e.equip_set_bind_seq(self.owner.GetSeq(), slot)
	} else {
		// swap
		box[slot].equip_set_bind_seq(0, slot)

		box[slot] = e
		e.equip_set_bind_seq(self.owner.GetSeq(), slot)
	}

	// ok
	return Err.OK
}

func (self *equipbox_t) __unequip(slot int32) int32 {
	box := self.box

	// check if slot is empty
	if box[slot] == nil {
		return Err.Equip_SlotIsEmpty
	}

	// do it
	box[slot].equip_set_bind_seq(0, slot)
	box[slot] = nil

	return Err.OK
}

func (self *equipbox_t) get_mods(grpids ...int) (ret PropMods) {
	ret = NewPropMods()

	if len(grpids) == 0 {
		for k := range equip_group_registry {
			grpids = append(grpids, k)
		}
	}

	for _, id := range grpids {
		grp := equip_group_registry[id]
		if grp != nil {
			grp.apply_mods_func(self.box, grp.slot_start, grp.slot_end, ret)
		}
	}

	return
}

// ============================================================================

func (self *EquipOp) Equip(e IEquippable, opt_slot ...int32) *EquipOp {
	// check if it's equipped
	if e.equip_get_bind_seq() > 0 {
		self.err = Err.Equip_AlreadyEquipped
		return self
	}

	// check slot
	var slot int32
	if len(opt_slot) > 0 {
		slot = opt_slot[0]
	} else {
		slot = e.equip_get_slot(self.ebox.box)
	}

	if !equip_is_slot_valid(e, slot) {
		self.err = Err.Equip_InvalidSlot
		return self
	}

	// add op
	self.ops = append(self.ops, &equip_op_entry_t{
		tp:   eop_Equip,
		e:    e,
		slot: slot,
	})

	return self
}

func (self *EquipOp) Unequip(slot int32) *EquipOp {
	// check slot
	if !equip_is_slot_valid(nil, slot) {
		self.err = Err.Equip_InvalidSlot
		return self
	}

	// check if slot is empty
	e := self.ebox.box[slot]
	if e == nil {
		self.err = Err.Equip_SlotIsEmpty
		return self
	}

	// add op
	self.ops = append(self.ops, &equip_op_entry_t{
		tp:   eop_Unequip,
		e:    e,
		slot: slot,
	})
	return self
}

func (self *EquipOp) UnequipAll(grpids ...int) *EquipOp {
	if len(grpids) == 0 {
		for k := range equip_group_registry {
			grpids = append(grpids, k)
		}
	}

	for _, id := range grpids {
		grp := equip_group_registry[id]
		if grp == nil {
			continue
		}

		for i := grp.slot_start; i <= grp.slot_end; i++ {
			e := self.ebox.box[i]
			if e != nil {
				self.ops = append(self.ops, &equip_op_entry_t{
					tp:   eop_Unequip,
					e:    e,
					slot: i,
				})
			}
		}
	}

	return self
}

func (self *EquipOp) Modify(e IEquippable, f func()) *EquipOp {
	// must be equipped
	if e.equip_get_bind_seq() != self.ebox.owner.GetSeq() {
		self.err = Err.Equip_NotEquipped
		return self
	}

	// add op
	self.ops = append(self.ops, &equip_op_entry_t{
		tp: eop_Modify,
		e:  e,
		f:  f,
	})

	return self
}

func (self *EquipOp) Apply() int32 {
	// check error
	if self.err != Err.OK {
		return self.err
	}

	// get affected groups
	grpids := self.affected_grps()

	// get old mods
	old_mods := self.ebox.get_mods(grpids...)

	// apply each op
	b := false
	for _, op := range self.ops {
		err := int32(Err.OK)

		switch op.tp {
		case eop_Equip:
			err = self.ebox.__equip(op.e, op.slot)
		case eop_Unequip:
			err = self.ebox.__unequip(op.slot)
		case eop_Modify:
			op.f()
		}

		if err == Err.OK {
			b = true
		} else {
			self.err = err
		}
	}

	// calc mods if at least some ops are executed successfully
	if b {
		// get new mods
		new_mods := self.ebox.get_mods(grpids...)

		// delta mods
		new_mods.Sub(old_mods)
		self.ebox.owner.GetMods().Add(new_mods)

		// calc props
		self.ebox.owner.CalcProps(true)
	}

	// return last error
	return self.err
}

func (self *EquipOp) affected_grps() (ret []int) {
outer:
	for _, op := range self.ops {
		grpid := op.e.equip_get_grpid()
		for _, v := range ret {
			if v == grpid {
				continue outer
			}
		}
		ret = append(ret, grpid)
	}

	return
}
