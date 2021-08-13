package worlddata

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/resetter"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
)

// ============================================================================

var g_resetter = &global_resetter_t{}

// ============================================================================

type global_resetter_t struct {
	resetter.Resettable
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_WorldReady, func(...interface{}) {
		resetter.Add(g_resetter)
	})
}

// ============================================================================

func (self *global_resetter_t) Open() {
	self.load_data()
}

func (self *global_resetter_t) load_data() {
	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_worlddata,
		"gresetter",
		&g_resetter,
	)
	if err != nil && !db.IsNotFound(err) {
		core.Panic("loading global resetter failed:", err)
	}

	log.Info("global resetter loaded:", self.Rst_ts)
}

func (self *global_resetter_t) Reset_Daily() {
	// fire
	evtmgr.Fire(gconst.Evt_GlobalResetDaily)

	// flush
	async.Push(func() {
		err := dbmgr.DBGame.Replace(
			dbmgr.C_tabname_worlddata,
			"gresetter",
			g_resetter,
		)
		if err != nil {
			log.Warning("flush global resetter failed:", err)
		}
	})
}

func (self *global_resetter_t) Reset_Weekly() {
	// fire
	evtmgr.Fire(gconst.Evt_GlobalResetWeekly)
}

func (self *global_resetter_t) Reset_Monthly() {
	// fire
	evtmgr.Fire(gconst.Evt_GlobalResetMonthly)
}
