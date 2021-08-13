package growfund

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mdata"
)

// ============================================================================

const NAME = "growfund"

func init() {
	mdata.Register(&mdata.Reg{
		Name: NAME,

		NewModuleData: new_data,
		DataLoaded:    data_loaded,
	})

	evtmgr.On(gconst.Evt_BillGrowFund, func(args ...interface{}) {
		GrowFundSvr.SvrCnt++
	})
}
