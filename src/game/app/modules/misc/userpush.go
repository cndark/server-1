package misc

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/shared/config"
	"net/url"
)

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_OfflineCounterFull, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		id := args[1].(int32)
		max := args[2].(int64)

		if id != gconst.Cnt_PlayerStrength || plr.IsOnline() {
			return
		}

		ar := plr.GetAuthRet()

		openid := ar["openid"]
		tpl := "ep_full"
		vals := []string{fmt.Sprintf("%d/%d", max, max)}

		wx_h5_push(openid, tpl, vals)
	})

	evtmgr.On(gconst.Evt_OfflineWlevelGjFull, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		if plr.IsOnline() {
			return
		}

		ar := plr.GetAuthRet()

		openid := ar["openid"]
		tpl := "gj_full"
		vals := []string{}

		wx_h5_push(openid, tpl, vals)
	})
}

// ============================================================================

func wx_h5_push(openid string, tpl string, vals []string) {
	if openid == "" {
		return
	}

	core.Go(func() {
		core.HttpPost(
			fmt.Sprintf(
				"http://%s:%d/api/userpush/wx_h5?token=%s",
				config.Reporter.IP,
				config.Reporter.Port,
				config.Reporter.Token,
			),
			url.Values{
				"openid": {openid},
				"tpl":    {tpl},
				"vals":   vals,
			},
		)
	})
}
