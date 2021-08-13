package calendar

import (
	"fw/src/core/evtmgr"
	"time"
)

// ============================================================================

func init() {
	evtmgr.On("cal->evt", func(args ...interface{}) {
		_ = args[0].(bool)
		name := args[1].(string)
		stage := args[2].(string)
		t0 := args[3].(time.Time)
		t1 := args[4].(time.Time)
		t2 := args[5].(time.Time)

		// handle stage
		handle_stage(name, stage, t0, t1, t2)
	})
}

func handle_stage(name, stage string, t0, t1, t2 time.Time) {
	// get reg
	reg := c_registry[name]
	if reg == nil {
		return
	}

	// get stage func
	sf := reg.StageFunc[stage]

	// on stage
	if reg.OnStage == nil {
		// trigger stage func
		if sf != nil {
			sf()
		}
	} else {
		reg.OnStage(stage, t0, t1, t2, func(b bool) {
			if b {
				// trigger stage func
				if sf != nil {
					sf()
				}
			}
		})
	}
}
