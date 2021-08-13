package warcup

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/calendar"
	"fw/src/game/app/modules/cond"
)

// ============================================================================

const (
	NAME = "warCup"
)

func init() {
	calendar.Register(&calendar.Reg{
		Name: NAME,

		OnStage: on_stage,

		StageFunc: map[string]func(){
			"open":     stage_open,
			"audition": stage_audition,
			"top64":    stage_top64,
			"top8":     stage_top8,
			"end":      stage_end,
			"close":    stage_close,
		},
	})

	evtmgr.On(gconst.Evt_GameDataLoaded, func(...interface{}) {
		for _, conf := range gamedata.ConfWarCupGuessTaskAttain.Items() {
			conf := conf

			cond.RegistCondObj(
				cond.Cond_Module_Plr_WarCup, conf.Id, conf.Cond, conf.P1,

				func(body interface{}) cond.ICondObj {
					return body.(IPlayer).GetWarCup().get(conf.Id)
				},
			)
		}
	})
}
