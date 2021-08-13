package gwar

import (
	"fw/src/game/app/modules/calendar"
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

const (
	NAME = "gwar"
)

func init() {
	mdata.Register(&mdata.Reg{
		Name:          NAME,
		NewModuleData: new_local_data,
		DataLoaded:    on_local_data_loaded,
	})

	calendar.Register(&calendar.Reg{
		Name: NAME,

		OnStage: on_stage,

		StageFunc: map[string]func(){
			"prepare": stage_prepare,
			"enroll":  stage_enroll,
			"match":   stage_match,
			"reward":  stage_reward,
			"close":   stage_close,
		},
	})
}
