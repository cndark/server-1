package act

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/worlddata"
	"time"
)

// ============================================================================

func init() {
	evtmgr.On("cal->act", func(args ...interface{}) {
		abs := args[0].(bool)
		name := args[1].(string)
		stage := args[2].(string)
		t0 := args[3].(time.Time)
		t1 := args[4].(time.Time)
		t2 := args[5].(time.Time)

		// handle stage
		handle_stage(abs, name, stage, t0, t1, t2)
	})
}

func handle_stage(abs bool, name string, stage string, t0, t1, t2 time.Time) {
	// find act
	act := FindAct(name)
	if act == nil {
		return
	}

	// check need-days for absolute time-line
	if abs {
		conf := gamedata.ConfActivity.Query(name)
		if conf != nil {
			need_ts := core.StartOfDay(worlddata.GetSvrCreateTs()).AddDate(0, 0, int(conf.NeedDays))
			if t0.Before(need_ts) {
				return
			}
		}
	}

	// on stage
	act.set_ver(t0)
	act.set_stage(stage)
	act.set_t1(t1)
	act.set_t2(t2)
	act.send_state_change()
	act.OnStage()
}
