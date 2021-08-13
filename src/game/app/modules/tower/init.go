package tower

import "fw/src/game/app/modules/mdata"

// ============================================================================

const NAME = "tower"

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME,

		NewModuleData: new_data,
		DataLoaded:    data_loaded,
	})
}
