package growfund

import (
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

var GrowFundSvr = &growfund_t{}

type growfund_t struct {
	SvrCnt int32
}

// ============================================================================

func new_data() interface{} {
	return &growfund_t{}
}

func data_loaded() {
	GrowFundSvr = mdata.Get(NAME).(*growfund_t)
}
