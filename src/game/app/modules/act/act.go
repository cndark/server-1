package act

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================

var (
	act_reg = make(map[string]IAct) // [act-name]
)

// ============================================================================

type act_rec_t struct {
	Ver     time.Time   // current ver
	Num     int32       // the nth number
	ActData *act_data_t // actdata
}

type act_data_t struct {
	Ver  time.Time   // data ver
	Data interface{} // real data
}

// ============================================================================

func Open() {
	// load act data
	load_act_data()

	// start save timer
	save_timer_start()

	// on init
	EachAct(func(a IAct) {
		a.OnInit()
	})
}

func Close() {
	// on close
	EachAct(func(a IAct) {
		a.OnQuit()
	})

	// save all acts
	save()
}

func RegisterAct(name string, a IAct) {
	if act_reg[name] != nil {
		log.Error("act already registered:", name)
		return
	}

	a.set_name(name)
	act_reg[name] = a
}

func FindAct(name string) IAct {
	a := act_reg[name]
	if a == nil {
		return nil
	} else {
		return a
	}
}

func EachAct(f func(a IAct)) {
	for _, a := range act_reg {
		f(a)
	}
}

// ============================================================================

func load_act_data() {
	// load
	var obj struct {
		Acts map[string]*act_rec_t
	}

	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_act,
		1,
		&obj,
	)
	if db.IsNotFound(err) {
		obj.Acts = make(map[string]*act_rec_t)
	} else if err != nil {
		core.Panic("loading act failed:", err)
	}

	// fill
	unix_t0 := time.Unix(0, 0)

	for _, a := range act_reg {
		rec := obj.Acts[a.GetName()]
		if rec == nil {
			// no rec. init it
			a.set_ver(unix_t0)
			a.set_num(0)
			a.set_stage("close")
			a.set_t2(unix_t0)
			a.set_actdata(&act_data_t{
				Ver:  unix_t0,
				Data: alloc_svr_data(nil, a),
			})
		} else {
			a.set_ver(rec.Ver)
			a.set_num(rec.Num)
			a.set_stage("close")
			a.set_t2(unix_t0)
			a.set_actdata(&act_data_t{
				Ver:  rec.ActData.Ver,
				Data: alloc_svr_data(rec.ActData.Data, a),
			})
		}
	}
}

func alloc_svr_data(m interface{}, a IAct) (out interface{}) {
	out = a.NewSvrData()

	if m == nil {
		return
	}

	data, err := bson.Marshal(m)
	if err != nil {
		log.Error("alloc act data failed:", err)
		return
	}

	err = bson.Unmarshal(data, out)
	if err != nil {
		log.Error("alloc act data failed:", err)
		return
	}

	return
}

// ============================================================================

func save_timer_start() {
	loop.SetTimeout(time.Now().Add(time.Duration(1200+rand.Intn(1200))*time.Second), func() {
		save_async()
		save_timer_start()
	})
}

// 异步存盘：用于定时存盘
func save_async() {
	// prepare data
	doc := save_gen_doc()

	// async save
	async.Push(func() {
		dbmgr.DBGame.Upsert(
			dbmgr.C_tabname_act,
			1,
			doc,
		)
	})
}

// 同步存盘：用于停服时存盘
func save() {
	doc := save_gen_doc()

	dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_act,
		1,
		doc,
	)
}

func save_gen_doc() db.M {
	acts := make(map[string]*act_rec_t)
	for _, a := range act_reg {
		acts[a.GetName()] = &act_rec_t{
			Ver:     a.get_ver(),
			Num:     a.GetNum(),
			ActData: a.get_actdata(),
		}
	}

	return db.M{
		"$set": db.M{"acts": core.CloneBsonObject(acts)},
	}
}
