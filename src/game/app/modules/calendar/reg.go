package calendar

import "time"

// ============================================================================

var (
	c_registry = make(map[string]*Reg) // [play-name]
)

// ============================================================================

type Reg struct {
	Name string

	OnStage   func(stg string, t0, t1, t2 time.Time, f func(bool)) // before each stage. f: pass false to prevent stage-func triggering
	StageFunc map[string]func()                                    // [stage]
}

// ============================================================================

func Register(reg *Reg) {
	c_registry[reg.Name] = reg
}
