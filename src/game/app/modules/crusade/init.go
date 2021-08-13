package crusade

import (
	"fw/src/game/app/modules/calendar"
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

const NAME = "crusade"

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME,

		NewModuleData: new_data,
		DataLoaded:    data_loaded,
	})

	calendar.Register(&calendar.Reg{
		Name:    NAME,
		OnStage: on_stage,
		StageFunc: map[string]func(){
			"start": stage_start,
			"end":   stage_end,
		},
	})
}
