package comp

// ============================================================================

// 属性Id 规则
//	* 固定点属性:     fixed = value                    [1,     99]
//	* 百分比加成属性: pct   = fixed * 100              [100,   9999]
//	* 聚集加成属性:   agr   = fixed * 10000 + count    [10000, 999999]
const (
	C_prop_min_pct_id       = 100
	C_prop_min_aggregate_id = 10000
)

// 固定属性 Id 含义
const ()

// ============================================================================
// PropMap
// ============================================================================

type PropMap map[int32]float32

// ============================================================================

func NewPropMap() PropMap {
	return make(PropMap)
}

func (self PropMap) Get(id int32) float32 {
	return self[id]
}

func (self PropMap) Set(id int32, v float32) {
	if v == 0 {
		delete(self, id)
	} else {
		self[id] = v
	}
}

func (self PropMap) Update(mods PropMods) (updated PropMap) {
	updated = NewPropMap()

	for id, p := range mods {
		if !p.dirty {
			continue
		}

		v := p.Base*(1+p.Pct) + p.Ext

		self.Set(id, v)
		updated[id] = v

		p.dirty = false
	}

	return
}

func (self *PropMap) Clear() {
	*self = NewPropMap()
}

// ============================================================================
// PropArray
// ============================================================================

type PropArray []*prop_t

type prop_t struct {
	Id  int32
	Val float32
}

// ============================================================================

func (self *PropArray) Append(id int32, v float32) {
	*self = append(*self, &prop_t{id, v})
}

func (self PropArray) Set(i int32, id int32, v float32) {
	if i < 0 || i >= int32(len(self)) {
		return
	}

	self[i].Id = id
	self[i].Val = v
}

// ============================================================================
// PropMods
// ============================================================================

type PropMods map[int32]*prop_entry_t

type prop_entry_t struct {
	Base  float32
	Pct   float32
	Ext   float32
	dirty bool
}

// ============================================================================

func NewPropMods() PropMods {
	return make(PropMods)
}

func (self PropMods) ModBase(id int32, v float32) {
	if v == 0 {
		return
	}

	if id < C_prop_min_pct_id {
		// normal
		p := self.find(id)
		p.Base += v
		p.dirty = true
	} else if id < C_prop_min_aggregate_id {
		// pct
		id /= C_prop_min_pct_id
		p := self.find(id)
		p.Pct += v
		p.dirty = true
	} else {
		// aggr
		fromid, count := id/C_prop_min_aggregate_id, id%C_prop_min_aggregate_id

		for id = fromid; id < fromid+count; id++ {
			p := self.find(id)
			p.Base += v
			p.dirty = true
		}
	}
}

func (self PropMods) ModExt(id int32, v float32) {
	if v == 0 {
		return
	}

	if id < C_prop_min_pct_id {
		// normal
		p := self.find(id)
		p.Ext += v
		p.dirty = true
	} else if id < C_prop_min_aggregate_id {
		// pct
		id /= C_prop_min_pct_id
		p := self.find(id)
		p.Pct += v
		p.dirty = true
	} else {
		// aggr
		fromid, count := id/C_prop_min_aggregate_id, id%C_prop_min_aggregate_id

		for id = fromid; id < fromid+count; id++ {
			p := self.find(id)
			p.Ext += v
			p.dirty = true
		}
	}
}

func (self PropMods) Add(m PropMods) {
	for id, p2 := range m {
		if p2.Base == 0 && p2.Pct == 0 && p2.Ext == 0 {
			continue
		}

		p1 := self.find(id)

		p1.Base += p2.Base
		p1.Pct += p2.Pct
		p1.Ext += p2.Ext

		p1.dirty = true
	}
}

func (self PropMods) Sub(m PropMods) {
	for id, p2 := range m {
		if p2.Base == 0 && p2.Pct == 0 && p2.Ext == 0 {
			continue
		}

		p1 := self.find(id)

		p1.Base -= p2.Base
		p1.Pct -= p2.Pct
		p1.Ext -= p2.Ext

		p1.dirty = true
	}
}

func (self *PropMods) Clear() {
	*self = NewPropMods()
}

// ============================================================================

func (self PropMods) find(id int32) *prop_entry_t {
	p := self[id]
	if p == nil {
		p = &prop_entry_t{}
		self[id] = p
	}

	return p
}
