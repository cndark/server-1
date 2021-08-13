package gm

import (
	"fw/src/shared/config"
	"net/http"
)

// ============================================================================

func handle_dev(req *http.Request) (r string, err error) {
	if !config.Common.DevMode {
		err = ErrNoKey
		return
	}

	plr, err := get_player(req)
	if err != nil {
		return
	}

	// set tut step
	tut_tp := get_string(req, "tut_tp")
	tut_step := get_int32(req, "tut_step")
	if tut_step >= 0 {
		plr.GetTutorial().Set(tut_tp, tut_step, "")
	}

	wlv_pass := get_int32(req, "wlv_pass")
	if wlv_pass > 0 && wlv_pass < 10000 {
		plr.GetWLevel().GMSetLevel(wlv_pass)
	}

	tower_pass := get_int32(req, "tower_pass")
	if tower_pass > 0 && wlv_pass < 10000 {
		plr.GetTower().GMSetLevel(tower_pass)
	}

	return
}
