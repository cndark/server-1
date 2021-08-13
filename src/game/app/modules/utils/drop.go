package utils

import (
	"fw/src/game/app/gamedata"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

var (
	rand_drop = rand.New(rand.NewSource(time.Now().Unix()))
)

// ============================================================================

type drop_t struct {
	Id int32
	N  int64
}

// ============================================================================

func Drop(plr IPlayer, id int32) (ret []*drop_t) {
	conf := gamedata.ConfDrop.Query(id)
	if conf == nil {
		return
	}

	// fixed drops
	for _, v := range conf.FixedDrop {
		ret = append(ret, &drop_t{
			Id: v.Id,
			N:  int64(v.N),
		})
	}

	// roll drops
	for _, v := range conf.RollDrop {
		// drop N times
		for i := int32(0); i < v.N; i++ {
			// check prob
			if rand_drop.Int31n(10000) < v.Prob {
				d := drop_one(v.Grp)
				if d != nil {
					ret = append(ret, d)
				}
			}
		}
	}

	// cond drops
	for _, v := range conf.CondDrop {
		// check cond
		d := drop_cond(plr, v)
		if len(d) > 0 {
			ret = append(ret, d...)
		}

	}

	return
}

func drop_one(grpid int32) (ret *drop_t) {
	arr, wsum := gamedata.ConfDropGrpM.Query(grpid)
	if arr == nil || wsum <= 0 {
		return
	}

	// select one
	p := rand_drop.Int31n(wsum) // [0, wsum)
	for _, v := range arr {
		p -= v.Weight

		if p < 0 {
			ret = &drop_t{
				Id: v.Id,
				N:  int64(v.Num),
			}
			break
		}
	}

	return
}

func drop_cond(plr IPlayer, id int32) (ret []*drop_t) {
	conf := gamedata.ConfCondDrop.Query(id)
	if conf == nil {
		return
	}

	if IsStatusTabArrayConform(plr, conf.Status) != Err.OK {
		return
	}

	for _, v := range conf.Items {
		p := rand_drop.Int31n(10000)
		if p <= v.Prob {
			ret = append(ret, &drop_t{
				Id: v.Id,
				N:  int64(v.N),
			})
		}
	}

	return
}
