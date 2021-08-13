package worlddata

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"fw/src/shared/config"
	"sync/atomic"
	"time"
)

// ============================================================================

const C_save_itv = 60 // in seconds

// ============================================================================

var (
	seq_hero  int64
	seq_armor int64
	seq_relic int64
)

// ============================================================================

type seq_t struct {
	Id  string `bson:"_id"`
	Seq int64  `bson:"seq"`
}

// ============================================================================

func seq_init() {
	seq_init_one("seq_hero", &seq_hero)
	seq_init_one("seq_armor", &seq_armor)
	seq_init_one("seq_relic", &seq_relic)

	seq_save_timer_start()
}

func seq_init_one(id string, seq *int64) {
	var obj seq_t

	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_worlddata,
		id,
		&obj,
	)
	if err == nil {
		*seq = obj.Seq
	} else if db.IsNotFound(err) {
		*seq = int64(config.CurGame.Id) << 40
	} else {
		core.Panic("init seq failed:", err, id)
	}
}

// ============================================================================

func seq_save_timer_start() {
	loop.SetTimeout(time.Now().Add(C_save_itv*time.Second), func() {
		seq_save_async()
		seq_save_timer_start()
	})
}

func seq_save_async() {
	async.Push(
		func() {
			seq_save_one("seq_hero", atomic.LoadInt64(&seq_hero))
			seq_save_one("seq_armor", atomic.LoadInt64(&seq_armor))
			seq_save_one("seq_relic", atomic.LoadInt64(&seq_relic))
		},
	)
}

func seq_save() {
	seq_save_one("seq_hero", seq_hero)
	seq_save_one("seq_armor", seq_armor)
	seq_save_one("seq_relic", seq_relic)
}

func seq_save_one(id string, seq int64) {
	err := dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_worlddata,
		id,
		db.M{"$set": db.M{"seq": seq}},
	)
	if err != nil {
		log.Warning("flush seq failed:", err, id)
	}
}

// ============================================================================

func GenSeqHero() int64 {
	return atomic.AddInt64(&seq_hero, 1)
}

func GenSeqArmor() int64 {
	return atomic.AddInt64(&seq_armor, 1)
}

func GenSeqRelic() int64 {
	return atomic.AddInt64(&seq_relic, 1)
}
