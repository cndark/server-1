package gsvc

import (
	"fw/src/core/log"
	"fw/src/game/app/gsvc/bill"
	"fw/src/game/app/gsvc/gm"
	"fw/src/shared/config"
	"net/http"
)

// ============================================================================

func Start() {
	go func() {
		// routers
		http.HandleFunc("/gm/service", gm.Dispatcher)
		http.HandleFunc("/bill/service", bill.Dispatcher)

		// start
		log.Notice("starting game service on:", config.CurGame.Svc)

		err := http.ListenAndServe(config.CurGame.Svc, nil)
		if err != nil {
			log.Error("game service exited:", err)
		}
	}()
}
