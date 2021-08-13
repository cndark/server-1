package monthticket

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/cond"
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GameDataLoaded, func(...interface{}) {
		for _, conf := range gamedata.ConfMonthTicketTask.Items() {
			conf := conf

			cond.RegistCondObj(
				cond.Cond_Module_Plr_MonthTicket, conf.Id, conf.Cond, conf.P1,

				func(body interface{}) cond.ICondObj {
					return body.(IPlayer).GetMonthTicket().get(conf.Id)
				},
			)
		}
	})
}
