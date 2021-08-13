package rift

import (
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

const NAME_MINE = "rift_mine"
const NAME_BOX = "rift_box"

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME_MINE,

		LoadData:   load_data_mine,
		DataLoaded: data_loaded_mine,

		SaveAsync: save_async_mine,
		Save:      save_mine,
	})

	mdata.Register(&mdata.Reg{
		Name: NAME_BOX,

		LoadData:   load_data_box,
		DataLoaded: data_loaded_box,

		SaveAsync: save_async_box,
		Save:      save_box,
	})
}
