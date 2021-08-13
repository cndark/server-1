package rift

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/game/app/dbmgr"
	"fw/src/game/msg"
	"math/rand"
	"time"
)

// ============================================================================

var rand_rift = rand.New(rand.NewSource(time.Now().Unix()))

const (
	C_RiftType_Mine    = 1 // 矿
	C_RiftType_Monster = 2 // 怪物
	C_RiftType_Box     = 3 // 宝箱
)

// 矿
var MineMgr = &rift_mine_t{
	Mines: make(map[int64]*mine_t),
	Plrs:  make(map[string]*plr_mine_t),
}

// 宝箱
var BoxMgr = &rift_box_t{}

// ============================================================================

func Rift_ToMsg(plr IPlayer) *msg.RiftData {
	ret := &msg.RiftData{}
	ret.Monster = plr.GetRift().ToMsg()
	ret.PlrMine = MineMgr.ToMsg(plr)
	ret.CurBox, ret.BoxNum = BoxMgr.ToMsg(plr)

	return ret
}

// ============================================================================
// 矿

func load_data_mine() interface{} {
	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_rift,
		"mine",
		&MineMgr,
	)

	if err != nil && !db.IsNotFound(err) {
		core.Panic("loading rift mine data failed:", err)
	}

	if MineMgr.Mines == nil {
		MineMgr.Mines = make(map[int64]*mine_t)
	}

	if MineMgr.Plrs == nil {
		MineMgr.Plrs = make(map[string]*plr_mine_t)
	}

	return MineMgr
}

func data_loaded_mine() {
	MineMgr.data_loaded()
}

func save_mine() {
	err := dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_rift,
		"mine",
		db.M{"$set": MineMgr},
	)

	if err != nil {
		log.Warning("save rift mine failed:", err)
	}

	log.Info("save rift mine")
}

func save_async_mine() {
	doc := core.CloneBsonObject(MineMgr)

	async.Push(func() {
		err := dbmgr.DBGame.Replace(
			dbmgr.C_tabname_rift,
			"mine",
			doc,
		)
		if err != nil {
			log.Error("save_async rift mine data failed:", err)
		}
	})
}

// ============================================================================
// 宝箱

func load_data_box() interface{} {
	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_rift,
		"box",
		&BoxMgr,
	)

	if err != nil && !db.IsNotFound(err) {
		core.Panic("loading rift box data failed:", err)
	}

	return BoxMgr
}

func data_loaded_box() {
	BoxMgr.data_loaded()
}

func save_box() {
	err := dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_rift,
		"box",
		db.M{"$set": BoxMgr},
	)

	if err != nil {
		log.Warning("save rift box failed:", err)
	}

	log.Info("save rift box")
}

func save_async_box() {
	doc := core.CloneBsonObject(BoxMgr)

	async.Push(func() {
		err := dbmgr.DBGame.Replace(
			dbmgr.C_tabname_rift,
			"box",
			doc,
		)
		if err != nil {
			log.Error("save_async rift box data failed:", err)
		}
	})
}
