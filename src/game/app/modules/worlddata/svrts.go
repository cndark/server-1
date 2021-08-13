package worlddata

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/game/app/dbmgr"
	"time"
)

// ============================================================================

var (
	svrts svrts_t
)

type svrts_t struct {
	Create_ts time.Time //
	Start_ts  time.Time //
	Merge_ts  time.Time // modified ONLY by merge script
	Merge_cnt int32     // modified ONLY by merge script
}

// ============================================================================

func svrts_init() {
	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_worlddata,
		"svrts",
		&svrts,
	)
	if err == nil {
		svrts_on_start(false)
	} else if db.IsNotFound(err) {
		svrts_on_start(true)
	} else {
		core.Panic("init svr time failed:", err)
	}
}

func GetSvrCreateTs() time.Time {
	return svrts.Create_ts
}

func GetSvrMergeTs() time.Time {
	return svrts.Merge_ts
}

func GetSvrMergeCnt() int32 {
	return svrts.Merge_cnt
}

func GetSvrStartTs() time.Time {
	return svrts.Start_ts
}

// ============================================================================

func svrts_on_start(isnew bool) {
	now := time.Now()

	// update ts
	doc := db.M{}

	if isnew {
		svrts.Create_ts = now
		svrts.Merge_ts = time.Unix(0, 0)

		doc["create_ts"] = svrts.Create_ts
		doc["merge_ts"] = svrts.Merge_ts
	}

	svrts.Start_ts = now
	doc["start_ts"] = svrts.Start_ts

	err := dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_worlddata,
		"svrts",
		db.M{"$set": doc},
	)
	if err != nil {
		log.Warning("flush svr ts failed:", err)
	}
}
