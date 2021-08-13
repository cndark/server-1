package app

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/shared/config"
	"net/url"
	"strconv"
	"time"
)

// ============================================================================

func init() {
	ch := make(chan bool)

	evtmgr.On(gconst.Evt_ServerStart, func(...interface{}) {
		core.Go(func() {
			defer func() { ch <- true }()
			for {
				select {
				case <-ch:
					return

				default:
					rpt_start()
					time.Sleep(time.Second * 10)
				}
			}
		})
	})

	evtmgr.On(gconst.Evt_ServerStop, func(...interface{}) {
		ch <- true
		<-ch
		close(ch)

		rpt_stop()
	})
}

func rpt_start() {
	core.HttpPost(
		fmt.Sprintf(
			"http://%s:%d/server/gate_add?token=%s",
			config.Switcher.IP,
			config.Switcher.Port,
			config.Switcher.Token,
		),
		url.Values{
			"name":   {config.CurGate.Name},
			"ip":     {config.CurGate.IPWan},
			"port":   {core.I32toa(config.CurGate.Port)},
			"wsport": {core.I32toa(config.CurGate.WsPort)},
			"load":   {strconv.FormatFloat(float64(NetMgr.session_count()*100)/C_max_session_count, 'f', 1, 64)},
		},
	)
}

func rpt_stop() {
	core.HttpPost(
		fmt.Sprintf(
			"http://%s:%d/server/gate_remove?token=%s",
			config.Switcher.IP,
			config.Switcher.Port,
			config.Switcher.Token,
		),
		url.Values{
			"name": {config.CurGate.Name},
		},
	)
}
