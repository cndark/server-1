package dbmgr

import (
	"fw/src/core/db"
	"fw/src/core/log"
)

// ============================================================================

func Share_UpdateUserName(userid string, name string) {
	err := DBShare.Update(
		C_tabname_userinfo,
		userid,
		db.M{"$set": db.M{"name": name}},
	)
	if err != nil {
		log.Error("dbmgr.Share_UpdateUserName() failed:", err)
	}
}

func Share_UpdateUserLv(userid string, lv int32) {
	err := DBShare.Update(
		C_tabname_userinfo,
		userid,
		db.M{"$set": db.M{"lv": lv}},
	)
	if err != nil {
		log.Error("dbmgr.Share_UpdateUserLv() failed:", err)
	}
}

func Share_UpdateUserVip(userid string, vip int32) {
	err := DBShare.Update(
		C_tabname_userinfo,
		userid,
		db.M{"$set": db.M{"vip": vip}},
	)
	if err != nil {
		log.Error("dbmgr.Share_UpdateUserVip() failed:", err)
	}
}
