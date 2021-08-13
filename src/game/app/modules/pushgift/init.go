package pushgift

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/cond"
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GameDataLoaded, func(...interface{}) {
		for _, conf := range gamedata.ConfPushGift.Items() {
			conf := conf

			cond.RegistCondObj(
				cond.Cond_Module_Plr_PushGift, conf.Id, conf.Cond, conf.P1,

				func(body interface{}) cond.ICondObj {
					if body.(IPlayer).GetPushGift().isfin(conf.Id) {
						return nil
					}

					var val val_t

					return &val
				},
			)
		}
	})
}
