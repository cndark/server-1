package perfmon

import (
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"os"
	"time"
)

// ============================================================================

const N = 5

// ============================================================================

func Start() {
	if os.Getenv("PERF_MON") != "true" {
		return
	}

	go func() {
		for {
			time.Sleep(N * time.Second)

			log.Infof("L:%5d, OL:%5d, LQ:%6d, H:%6d, AQ:%6d,%6d",
				app.PlayerMgr.NumLoaded(),
				app.PlayerMgr.NumOnline(),
				loop.QLen(),
				loop.NumHandled()/N,
				async.QLen(gconst.AQ_Default),
				async.QLen(gconst.AQ_GLog),
			)
		}
	}()
}
