package robot

import "fw/src/game/app/modules/mdata"

// ============================================================================

const NAME = "robot"

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME,

		NewModuleData: new_data,
		DataLoaded:    data_loaded,
	})
}
