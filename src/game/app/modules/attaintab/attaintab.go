package attaintab

import (
	"fw/src/game/msg"
)

// ============================================================================

type AttainTab struct {
	Objs map[int32]*attain_obj_t // [id]item

	plr IPlayer
}

type attain_obj_t struct {
	Id  int32
	Val float64 // progress value
}

// ============================================================================

func NewAttainTab() *AttainTab {
	return &AttainTab{
		Objs: make(map[int32]*attain_obj_t),
	}
}

// ============================================================================

func (self *AttainTab) Init(plr IPlayer) {
	self.plr = plr
}

func (self *AttainTab) get(oid int32) *attain_obj_t {
	obj := self.Objs[oid]
	if obj == nil {
		self.Objs[oid] = &attain_obj_t{Id: oid}
		return self.Objs[oid]
	}

	return obj
}

func (self *AttainTab) GetObjVal(oid int32) float64 {
	obj := self.get(oid)

	return obj.Val
}

// ============================================================================
func (self *AttainTab) ToMsg() *msg.AttainTabData {
	ret := &msg.AttainTabData{}

	for _, v := range self.Objs {
		ret.Objs = append(ret.Objs, &msg.AttainObj{
			OId: v.Id,
			Val: v.Val,
		})
	}

	return ret
}

// ============================================================================
// implements ICondObj interface

func (self *attain_obj_t) GetVal() float64 {
	return self.Val
}

func (self *attain_obj_t) SetVal(v float64) {
	self.Val = v
}

func (self *attain_obj_t) AddVal(v float64) {
	self.Val += v
}

func (self *attain_obj_t) Done(body interface{}, confid int32, isChange bool) {
	if !isChange {
		return
	}

	body.(IPlayer).SendMsg(&msg.GS_AttainTabObjValueChanged{
		OId: self.Id,
		Val: self.Val,
	})
}

// ============================================================================
