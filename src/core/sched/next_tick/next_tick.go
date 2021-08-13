package next_tick

import (
	"fw/src/core"
	"fw/src/core/sched/loop"
	"time"
)

// ============================================================================

var (
	f_arr []func()
	f_map = map[string]func(){}

	tid *core.Timer
)

// ============================================================================

func Push(f func()) {
	// add
	f_arr = append(f_arr, f)

	// check timer
	if tid == nil {
		tid = loop.SetTimeout(time.Now(), on_tick)
	}
}

func Once(key string, f func()) {
	// check key
	if f_map[key] != nil {
		return
	}

	// add
	if f_map == nil {
		f_map = make(map[string]func())
	}
	f_map[key] = f

	// check timer
	if tid == nil {
		tid = loop.SetTimeout(time.Now(), on_tick)
	}
}

// ============================================================================

func on_tick() {
	for i, f := range f_arr {
		f_arr[i] = nil
		f()
	}

	for _, f := range f_map {
		f()
	}

	f_arr = f_arr[:0]
	f_map = nil

	tid = nil
}
