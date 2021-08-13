package mopen

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/cond"
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GameDataLoaded, func(...interface{}) {
		for _, conf := range gamedata.ConfOpen.Items() {
			conf := conf

			cond.RegistCondObj(
				cond.Cond_Module_Plr_MOpen, conf.ModuleId, conf.Cond, conf.P1,

				func(body interface{}) cond.ICondObj {
					if body.(IPlayer).GetMOpen().IsOpen(conf.ModuleId) {
						return nil
					}

					return body.(IPlayer).GetMOpen().get(conf.ModuleId)
				},
			)
		}
	})
}
