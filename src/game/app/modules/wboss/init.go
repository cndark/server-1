package wboss

import (
	"fw/src/game/app/modules/calendar"
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

const (
	NAME = "wboss"
)

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME,

		NewModuleData: new_svrdata,
		DataLoaded:    on_svrdata_loaded,
	})

	calendar.Register(&calendar.Reg{
		Name: NAME,

		OnStage: on_stage,

		StageFunc: map[string]func(){
			"open":   stage_open,
			"start":  stage_start,
			"end":    stage_end,
			"reward": stage_reward,
			"close":  stage_close,
		},
	})
}
