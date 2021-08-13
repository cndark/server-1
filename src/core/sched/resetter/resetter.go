package resetter

import (
	"fw/src/core"
	"fw/src/core/sched/loop"
	"time"
)

// ============================================================================

type Resettable struct {
	Rst_ts time.Time // last reset timestamp
}

func (self *Resettable) Reset_GetTime() time.Time {
	return self.Rst_ts
}

func (self *Resettable) Reset_SetTime(ts time.Time) {
	self.Rst_ts = ts
}

type IResettable interface {
	Reset_GetTime() time.Time
	Reset_SetTime(ts time.Time)
	Reset_Daily()
	Reset_Weekly()
	Reset_Monthly()
}

// ============================================================================

var (
	objs = make(map[IResettable]bool) // obj map
)

// ============================================================================

func Start() {
	// sched for next reset
	key := calc_keytime()
	sched(key)
}

func Add(obj IResettable) {
	objs[obj] = true

	// check if we should reset
	key := calc_keytime()
	check_reset(obj, key)
}

func Remove(obj IResettable) {
	delete(objs, obj)
}

// ============================================================================

// return lastest key time before now
func calc_keytime() time.Time {
	now := time.Now()

	y, M, d := now.Date()
	key := time.Date(y, M, d, 0, 0, 0, 0, time.Local)

	if key.After(now) {
		key = key.Add(-24 * time.Hour)
	}

	return key
}

func sched(key time.Time) {
	key = key.Add(24 * time.Hour)
	loop.SetTimeout(key, func() {
		check_objs(key)
		sched(key)
	})
}

func check_objs(key time.Time) {
	for obj, _ := range objs {
		check_reset(obj, key)
	}
}

func check_reset(obj IResettable, key time.Time) {
	// check
	last_ts := obj.Reset_GetTime()
	if !last_ts.Before(key) {
		return
	}

	// update reset time
	obj.Reset_SetTime(key)

	// daily
	obj.Reset_Daily()

	// weekly
	if !core.IsSameWeek(last_ts, key) {
		obj.Reset_Weekly()
	}

	// monthly
	if !core.IsSameMonth(last_ts, key) {
		obj.Reset_Monthly()
	}
}
