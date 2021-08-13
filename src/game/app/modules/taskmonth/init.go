package taskmonth

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/cond"
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GameDataLoaded, func(...interface{}) {
		for _, conf := range gamedata.ConfTaskMonth.Items() {
			conf := conf

			cond.RegistCondObj(
				cond.Cond_Module_Plr_TaskMonth, conf.Id, conf.Cond, conf.P1,

				func(body interface{}) cond.ICondObj {
					return body.(IPlayer).GetTaskMonth().get(conf.Id)
				},
			)
		}
	})

	evtmgr.On(gconst.Evt_PlrResetMonthly, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr_data := plr.GetTaskMonth()
		if plr_data == nil {
			return
		}

		plr_data.reset_month()
	})
}
