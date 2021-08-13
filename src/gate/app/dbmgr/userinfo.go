package dbmgr

import (
	"fmt"
	"fw/src/core/db"
	"fw/src/core/log"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

const (
	C_player_name_maxlen = 18
)

// ============================================================================

type UserInfo struct {
	UserId   string    `bson:"_id"`
	AuthId   string    `bson:"authid"`
	Area     int32     `bson:"area"`
	Svr0     string    `bson:"svr0"`
	Svr      string    `bson:"svr"`
	Sdk      string    `bson:"sdk"`
	Model    string    `bson:"model"`
	DevId    string    `bson:"devid"`
	Os       string    `bson:"os"`
	OsVer    string    `bson:"osver"`
	IP       string    `bson:"ip"`
	CreateTs time.Time `bson:"cts"`
	BanTs    time.Time `bson:"ban_ts"`
	Name     string    `bson:"name"`
	Lv       int32     `bson:"lv"`
	Vip      int32     `bson:"vip"`
}

// ============================================================================

func Center_GetMergedSvrInfo(svr0 string) (svr string, closereg bool) {
	var obj struct {
		MName    string
		CloseReg bool
	}

	// get merged info
	err := DBCenter.GetObject(
		C_tabname_svrlist,
		svr0,
		&obj,
	)
	if err != nil {
		return
	}

	svr = obj.MName
	closereg = obj.CloseReg

	// correct closereg
	if svr != svr0 {
		err = DBCenter.GetObject(
			C_tabname_svrlist,
			svr,
			&obj,
		)
		if err != nil {
			return "", false
		}

		closereg = obj.CloseReg
	}

	// ok
	return
}

func Center_UpdateAcctLastSvr(auth_id string, sdk string, svr0 string) {
	err := DBCenter.Upsert(
		C_tabname_acctinfo,
		fmt.Sprintf("%s-%s", sdk, auth_id),
		db.M{
			"$set": db.M{
				"lastsvr": svr0,
			},
		},
	)
	if err != nil {
		log.Warning("update acct last svr failed:", err)
	}
}

func Center_UpdateAcctAuthRet(auth_id string, sdk string, authRet map[string]string) {
	if len(authRet) == 0 {
		return
	}

	err := DBCenter.Upsert(
		C_tabname_acctinfo,
		fmt.Sprintf("%s-%s", sdk, auth_id),
		db.M{
			"$set": db.M{
				"authret": authRet,
			},
		},
	)
	if err != nil {
		log.Warning("update acct authret failed:", err)
	}
}

func Center_GetUserInfo(auth_id string, sdk string, svr0, svr string, closereg bool, on_create func(*UserInfo)) (int32, *UserInfo) {
	var obj UserInfo

	err := DBShare.GetObjectByCond(
		C_tabname_userinfo,
		db.M{
			"authid": auth_id,
			"sdk":    sdk,
			"area":   config.Common.Area.Id,
			"svr0":   svr0,
		},
		&obj,
	)
	if err == nil {
		// userinfo exists
		// check svr
		if obj.Svr == svr {
			return Err.OK, &obj
		} else {
			return Err.Login_UserSvr, nil
		}
	} else if db.IsNotFound(err) {
		// check closereg
		if closereg {
			return Err.Login_CloseReg, nil
		}

		// allocate user db
		dbname := Center_AllocUserDB()
		if dbname == "" {
			return Err.Login_Failed, nil
		}

		// on create
		if on_create != nil {
			on_create(&obj)
		}

		// create new userinfo
		obj.UserId = Center_GenUserId(dbname)
		obj.AuthId = auth_id
		obj.Area = config.Common.Area.Id
		obj.Svr0 = svr0
		obj.Svr = svr
		obj.Sdk = sdk
		obj.CreateTs = time.Now()
		obj.BanTs = time.Unix(0, 0)
		obj.Lv = 1
		obj.Vip = 0

		// flush to db
		err = DBShare.Insert(C_tabname_userinfo, &obj)
		if err == nil {
			// update user load
			Center_IncUserLoad(dbname)

			// return new userinfo
			return Err.OK, &obj
		} else {
			// failed
			return Err.Login_CreateUserInfo, nil
		}
	} else {
		// failed
		return Err.Failed, nil
	}
}
