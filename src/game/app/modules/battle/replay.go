package battle

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"sync/atomic"
	"time"
)

// ============================================================================

var (
	seq_replay_id int32
)

// ============================================================================

func ReplaySave(xdb *db.Database, coll string, rp *msg.BattleReplay) string {
	now := time.Now().Unix()
	rp.Id = fmt.Sprintf(
		"%d-%d-%d",
		config.CurGame.Id, // svrid
		now,               // time
		atomic.AddInt32(&seq_replay_id, 1)%100000, // seq
	)
	rp.Ts = now

	async.Push(func() {
		err := xdb.Insert(coll, db.M{
			"_id": rp.Id,
			"d":   rp,
		})
		if err != nil {
			log.Error("saving battle replay failed:", rp, err)
		}
	})

	return rp.Id
}

func ReplayGet(xdb *db.Database, coll string, id string, f func(*msg.BattleReplay)) {
	core.Go(func() {
		var rec struct {
			D *msg.BattleReplay
		}
		err := xdb.GetObject(coll, id, &rec)

		loop.Push(func() {
			if err != nil {
				f(nil)
			} else {
				f(rec.D)
			}
		})
	})
}

func ReplayDel(xdb *db.Database, coll string, id string) {
	core.Go(func() {
		err := xdb.Remove(coll, id)
		if err != nil {
			log.Error("del battle replay failed:", coll, err)
		}
	})
}
