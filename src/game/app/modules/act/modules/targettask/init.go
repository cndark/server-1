package targettask

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/cond"
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GameDataLoaded, func(...interface{}) {
		for _, conf := range gamedata.ConfActTargetTaskAttain.Items() {
			conf := conf
			cond.RegistCondObj(cond.Cond_Module_Plr_ActTargetTask, conf.AttainId, conf.Cond, conf.P1,

				func(body interface{}) cond.ICondObj {
					if !actObj.Started() {
						return nil
					}

					plr_data := actObj.GetPlrData(body.(IPlayer))
					if plr_data == nil {
						return nil
					}

					return plr_data.get_attain_obj(conf.AttainId)
				})
		}
	})
}
