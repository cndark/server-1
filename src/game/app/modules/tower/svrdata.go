package tower

import (
	"fw/src/game/app/modules/mdata"
	"fw/src/game/msg"
)

// ============================================================================

var (
	TowerMgr = &tower_mgr_t{
		Records: make(record_m),
	}
)

// ============================================================================

type tower_mgr_t struct {
	Records record_m
}

type record_m map[int32]*record_t
type record_t struct {
	First    *msg.BattleReplay // 首杀战斗录像
	MinPower *msg.BattleReplay // 最小战力战斗录像
}

// ============================================================================

func new_data() interface{} {
	return &tower_mgr_t{
		Records: make(record_m),
	}
}

func data_loaded() {
	TowerMgr = mdata.Get(NAME).(*tower_mgr_t)
}

// ============================================================================

func (self *tower_mgr_t) add_replay(lv int32, r *msg.BattleReplay) {
	rec := self.Records[lv]
	if rec == nil {
		rec = &record_t{
			First:    r,
			MinPower: r,
		}

		self.Records[lv] = rec
		return
	}

	if r == nil || r.Bi == nil || r.Bi.T1 == nil || r.Bi.T1.Fighters == nil {
		return
	}

	if rec.MinPower == nil || rec.MinPower.Bi == nil || rec.MinPower.Bi.T1 == nil ||
		rec.MinPower.Bi.T1.Fighters == nil {
		rec.MinPower = r
		return
	}

	min := int32(0)
	for _, v := range rec.MinPower.Bi.T1.Fighters {
		min += v.AtkPwr
	}

	cur := int32(0)
	for _, v := range r.Bi.T1.Fighters {
		cur += v.AtkPwr
	}

	if cur < min {
		rec.MinPower = r
	}
}
