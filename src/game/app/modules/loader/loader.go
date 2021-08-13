package loader

import (
	_ "fw/src/game/app/modules/act/modules/actgift"
	_ "fw/src/game/app/modules/act/modules/billltday"
	_ "fw/src/game/app/modules/act/modules/billlttotal"
	_ "fw/src/game/app/modules/act/modules/heroskin"
	_ "fw/src/game/app/modules/act/modules/monopoly"
	_ "fw/src/game/app/modules/act/modules/msummon"
	_ "fw/src/game/app/modules/act/modules/rushlocal"
	_ "fw/src/game/app/modules/act/modules/summon"
	_ "fw/src/game/app/modules/act/modules/targettask"
	_ "fw/src/game/app/modules/arena"
	_ "fw/src/game/app/modules/attaintab"
	_ "fw/src/game/app/modules/crusade"
	_ "fw/src/game/app/modules/growfund"
	_ "fw/src/game/app/modules/gwar"
	_ "fw/src/game/app/modules/ladder"
	_ "fw/src/game/app/modules/lamp"
	_ "fw/src/game/app/modules/monthticket"
	_ "fw/src/game/app/modules/mopen"
	_ "fw/src/game/app/modules/pushgift"
	_ "fw/src/game/app/modules/rift"
	_ "fw/src/game/app/modules/robot"
	_ "fw/src/game/app/modules/taskdaily"
	_ "fw/src/game/app/modules/taskmonth"
	_ "fw/src/game/app/modules/tower"
	_ "fw/src/game/app/modules/warcup"
	_ "fw/src/game/app/modules/wboss"
)

// ============================================================================

func LoadModules() {
	// for those listener modules that do NOT load by themselves
	// just import the modules here for loading

}
