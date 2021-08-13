package app

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/shared/config"
	"net/url"
	"sync/atomic"
	"time"
)

// ============================================================================

var (
	rpt_quit = make(chan int)

	rpt_new_ucnt      int32
	rpt_new_ucnt_last int32
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_ServerStart, func(...interface{}) {
		rpt_start()
		rpt_run_ticker()
	})

	evtmgr.On(gconst.Evt_ServerStop, func(...interface{}) {
		rpt_quit <- 1
		<-rpt_quit
		close(rpt_quit)

		rpt_new_user()
		rpt_stop()
	})

	evtmgr.On(gconst.Evt_ServerNewUser, func(args ...interface{}) {
		n := args[0].(int32)

		atomic.StoreInt32(&rpt_new_ucnt, n)
	})
}

func rpt_run_ticker() {
	ticker := time.NewTicker(time.Second * 10)

	core.Go(func() {
		defer func() {
			ticker.Stop()
			rpt_quit <- 1
		}()

		for {
			select {
			case <-rpt_quit:
				return

			case <-ticker.C:
				rpt_start()
				rpt_new_user()
			}
		}
	})
}

func rpt_start() {
	core.HttpPost(
		fmt.Sprintf(
			"http://%s:%d/server/game_add?token=%s",
			config.Switcher.IP,
			config.Switcher.Port,
			config.Switcher.Token,
		),
		url.Values{
			"name": {config.CurGame.Name},
		},
	)
}

func rpt_stop() {
	core.HttpPost(
		fmt.Sprintf(
			"http://%s:%d/server/game_remove?token=%s",
			config.Switcher.IP,
			config.Switcher.Port,
			config.Switcher.Token,
		),
		url.Values{
			"name": {config.CurGame.Name},
		},
	)
}

func rpt_new_user() {
	n := atomic.LoadInt32(&rpt_new_ucnt)
	if n == rpt_new_ucnt_last {
		return
	}
	rpt_new_ucnt_last = n

	core.HttpPost(
		fmt.Sprintf(
			"http://%s:%d/newucnt",
			config.Agent.IP,
			config.Agent.Port,
		),
		url.Values{
			"area": {core.I32toa(config.Common.Area.Id)},
			"name": {config.CurGame.Name},
			"n":    {core.I32toa(n)},
		},
	)
}
