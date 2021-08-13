package dbmgr

import (
	"fw/src/core/db"
	"fw/src/shared/config"
)

// ============================================================================

const (
	// share
	C_tabname_userinfo = "userinfo"

	// center
	C_tabname_acctinfo = "acctinfo"
	C_tabname_svrlist  = "svrlist"
	C_tabname_userload = "userload"
	C_tabname_seqid    = "seqid"
	C_tabname_names    = "names"
)

// ============================================================================

var (
	DBShare  *db.Database
	DBCenter *db.Database
)

// ============================================================================

func Open() {
	open_share()
	open_center()
}

func Close() {
	DBShare.Close()
	DBCenter.Close()
}

// ============================================================================

func open_share() {
	if DBShare != nil {
		return
	}

	DBShare = db.NewDatabase()
	DBShare.Open(config.Common.DBShare, 1)
}

func open_center() {
	if DBCenter != nil {
		return
	}

	DBCenter = db.NewDatabase()
	DBCenter.Open(config.Common.DBCenter, 1)
}
