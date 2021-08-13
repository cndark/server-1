package dbmgr

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/kfk"
	"fw/src/shared/config"

	"github.com/mediocregopher/radix/v3"
)

// ============================================================================

const (
	// share
	C_tabname_userinfo = "userinfo"
	C_tabname_giftinfo = "giftinfo"
	C_tabname_giftuse  = "giftuse"

	// center
	C_tabname_acctinfo = "acctinfo"
	C_tabname_userload = "userload"
	C_tabname_seqid    = "seqid"
	C_tabname_names    = "names"
	C_tabname_gnames   = "gnames"

	// stats
	C_tabname_login    = "login"
	C_tabname_bill     = "bill"
	C_tabname_online   = "online"
	C_tabname_tutorial = "tutorial"
	C_tabname_wlevel   = "wlevel"

	// log
	C_tabname_log  = "log"
	C_tabname_chat = "chat"

	// cross
	C_tabname_rank   = "rank"
	C_tabname_gwar   = "gwar"
	C_tabname_ladder = "ladder"

	// user
	C_tabname_user  = "user"
	C_tabname_guild = "guild"

	// game
	C_tabname_worlddata     = "worlddata"
	C_tabname_gmail         = "gmail"
	C_tabname_act           = "act"
	C_tabname_mdata         = "mdata"
	C_tabname_replayarena   = "replayarena"
	C_tabname_rift          = "rift"
	C_tabname_replay_ladder = "replay_ladder"
)

// ============================================================================

var (
	DBShare  *db.Database
	DBCenter *db.Database
	DBStats  *db.Database
	DBLog    *db.Database
	DBBill   *db.Database
	DBCross  *db.Database
	DBGame   *db.Database

	Kfk   *kfk.Kafka
	Redis *radix.Pool
)

var db_user = map[string]*db.Database{}

// ============================================================================

func Open() {
	open_share()
	open_center()
	open_stats()
	open_log()
	open_bill()
	open_cross()
	open_user()
	open_game()

	open_kfk()
	open_redis()
}

func Close() {
	DBShare.Close()
	DBCenter.Close()
	DBStats.Close()
	DBLog.Close()
	DBBill.Close()
	DBCross.Close()

	for _, db := range db_user {
		db.Close()
	}

	DBGame.Close()

	if Kfk != nil {
		Kfk.Close()
	}

	if Redis != nil {
		Redis.Close()
	}
}

func UserDB(dbname string) *db.Database {
	return db_user[dbname]
}

// ============================================================================

func open_share() {
	if DBShare != nil {
		return
	}

	DBShare = db.NewDatabase()
	DBShare.Open(config.Common.DBShare, 1)

	DBShare.CreateIndex(C_tabname_userinfo, "uk_authid", []string{"authid", "sdk", "area", "svr0"}, true)
	DBShare.CreateIndex(C_tabname_userinfo, "idx_area", []string{"area"}, false)
	DBShare.CreateIndex(C_tabname_userinfo, "idx_svr", []string{"svr"}, false)
	DBShare.CreateIndex(C_tabname_userinfo, "idx_sdk", []string{"sdk"}, false)
	DBShare.CreateIndex(C_tabname_userinfo, "idx_name", []string{"name"}, false)
	DBShare.CreateIndex(C_tabname_userinfo, "idx_devid", []string{"devid"}, false)
}

func open_center() {
	if DBCenter != nil {
		return
	}

	DBCenter = db.NewDatabase()
	DBCenter.Open(config.Common.DBCenter, 1)

	Center_CreateSeqId()
	Center_CreateUserLoad()

	DBCenter.CreateIndex(C_tabname_names, "uk_name", []string{"name"}, true)
	DBCenter.CreateIndex(C_tabname_gnames, "uk_gname", []string{"gname"}, true)
}

func open_stats() {
	if DBStats != nil {
		return
	}

	DBStats = db.NewDatabase()
	DBStats.Open(config.Common.DBStats, 1)

	DBStats.CreateIndex(C_tabname_login, "idx_area", []string{"area"}, false)
	DBStats.CreateIndex(C_tabname_login, "idx_svr", []string{"svr"}, false)
	DBStats.CreateIndex(C_tabname_login, "idx_sdk", []string{"sdk"}, false)
	DBStats.CreateIndex(C_tabname_login, "idx_uid", []string{"uid"}, false)
	DBStats.CreateIndex(C_tabname_login, "idx_cts", []string{"cts"}, false)
	DBStats.CreateIndex(C_tabname_login, "idx_lts", []string{"lts"}, false)
	DBStats.CreateIndex(C_tabname_login, "idx_day", []string{"day"}, false)

	DBStats.CreateIndex(C_tabname_bill, "idx_area", []string{"area"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_svr", []string{"svr"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_sdk", []string{"sdk"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_uid", []string{"uid"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_cts", []string{"cts"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_bts", []string{"bts"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_day", []string{"day"}, false)
	DBStats.CreateIndex(C_tabname_bill, "idx_amt", []string{"amt"}, false)

	DBStats.CreateIndex(C_tabname_online, "idx_area", []string{"area"}, false)
	DBStats.CreateIndex(C_tabname_online, "idx_svr", []string{"svr"}, false)
	DBStats.CreateIndex(C_tabname_online, "idx_ts", []string{"ts"}, false)
}

func open_log() {
	if DBLog != nil {
		return
	}

	DBLog = db.NewDatabase()
	DBLog.Open(config.Common.DBLog, 2)

	DBLog.CreateIndex(C_tabname_log, "idx_op", []string{"op"}, false)
	DBLog.CreateIndex(C_tabname_log, "idx_area", []string{"area"}, false)
	DBLog.CreateIndex(C_tabname_log, "idx_uid", []string{"uid"}, false)
	DBLog.CreateIndex(C_tabname_log, "idx_svr", []string{"svr"}, false)
	DBLog.CreateIndex(C_tabname_log, "idx_sdk", []string{"sdk"}, false)
	DBLog.CreateIndex(C_tabname_log, "idx_ts", []string{"ts"}, false)
}

func open_bill() {
	if DBBill != nil {
		return
	}

	DBBill = db.NewDatabase()
	DBBill.Open(config.Common.DBBill, 1)
}

func open_cross() {
	if DBCross != nil {
		return
	}

	DBCross = db.NewDatabase()
	DBCross.Open(config.Common.DBCross, 1)

	DBCross.CreateIndex(C_tabname_rank, "idx_type", []string{"type"}, false)
	DBCross.CreateIndex(C_tabname_rank, "idx_sgid", []string{"sgid"}, false)
	DBCross.CreateIndex(C_tabname_rank, "idx_rankid", []string{"rankid"}, false)
}

func open_user() {
	for k, v := range config.Common.DBUser {
		if db_user[k] != nil {
			continue
		}

		db := db.NewDatabase()
		db.Open(v, 1)

		db_user[k] = db

		db.CreateIndex(C_tabname_user, "idx_svr", []string{"base.svr"}, false)
		db.CreateIndex(C_tabname_user, "idx_lts", []string{"base.login_ts"}, false)

		db.CreateIndex(C_tabname_guild, "idx_svr", []string{"base.svr"}, false)
	}
}

func open_game() {
	if DBGame != nil {
		return
	}

	DBGame = db.NewDatabase()
	DBGame.Open(config.CurGame.DBGame, 1)
}

func open_kfk() {
	if Kfk != nil || len(config.Common.Kfk.Urls) == 0 {
		return
	}

	Kfk = kfk.NewKafka()
	Kfk.Open(config.Common.Kfk.Urls)
}

func open_redis() {
	if Redis != nil || len(config.Common.Redis.Urls) == 0 {
		return
	}

	p, err := radix.NewPool("tcp", config.Common.Redis.Urls[0], config.Common.Redis.PoolSize)
	if err != nil {
		core.Panic("create redis pool failed:", err)
	}

	err = p.Do(radix.Cmd(nil, "PING"))
	if err != nil {
		core.Panic("redis PING failed:", err)
	}

	Redis = p
}
