package __tpl

import (
	"fw/src/core"
	"fw/src/core/log"
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

var (
	m_data *moduledata_t
)

// ============================================================================

const NAME = "__tpl"

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME,

		NewModuleData: new_moduledata,

		LoadData:   load_data,
		DataLoaded: data_loaded,

		SaveAsync: save_async,
		Save:      save,
	})
}

// ============================================================================

type moduledata_t struct {
	Door string
	Num  int32
}

// ============================================================================

func new_moduledata() interface{} {
	return &moduledata_t{}
}

func load_data() interface{} {
	core.Panic("loading data failed")
	return nil
}

func data_loaded() {
	log.Warning("you can do initialization here")
}

func save_async() {
}

func save() {
}

// ============================================================================

func test() {
	log.Warning(m_data.Door, m_data.Num)
}
