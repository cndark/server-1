package battle

import (
	"fmt"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/shared/config"
	"math/rand"
)

// ============================================================================

var (
	// load-balancer for bats
	bat_arr   []*bat_t
	bat_arr_L int
	bat_arr_i int
)

// ============================================================================

type bat_t struct {
	addr string
}

// ============================================================================

func Init() {
	init_bats()

	evtmgr.On(gconst.Evt_ConfReload, func(args ...interface{}) {
		init_bats()
	})
}

// ============================================================================

func init_bats() {
	bat_arr = bat_arr[:0]
	for _, v := range config.Bats {
		bat_arr = append(bat_arr, &bat_t{
			addr: fmt.Sprintf("http://%s:%d/api/calc", v.IP, v.Port),
		})
	}
	bat_arr_L = len(bat_arr)

	rand.Shuffle(bat_arr_L, func(i, j int) {
		bat_arr[i], bat_arr[j] = bat_arr[j], bat_arr[i]
	})

	bat_arr_i = 0
}

func choose_bat() *bat_t {
	if bat_arr_L == 0 {
		return nil
	}

	bat_arr_i = (bat_arr_i + 1) % bat_arr_L

	return bat_arr[bat_arr_i]
}
