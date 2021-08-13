package misc

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/game/app/dbmgr"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

const (
	C_gift_code_len  = 8 + 4 + 4
	C_gift_grpid_len = 4
)

// ============================================================================

type GiftReward struct {
	ResK string `bson:"res_k"`
	ResV int64  `bson:"res_v"`
}

type GiftInfo struct {
	GrpId   int32         `bson:"_id"`
	Area    int32         `bson:"area"`
	Reuse   int32         `bson:"reuse"`
	Expire  time.Time     `bson:"expire"`
	Rewards []*GiftReward `bson:"rewards"`
}

// ============================================================================

func gift_extract_grpid(code string) int32 {
	if len(code) != C_gift_code_len {
		return 0
	}

	a := make([]byte, 0, C_gift_grpid_len)
	for i := 0; i < C_gift_grpid_len; i++ {
		c := code[i*3+2]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
			a = append(a, c)
		} else {
			break
		}
	}

	return core.A16toi32(string(a))
}

func gift_get_info(code string) *GiftInfo {
	// extract grpid
	grpid := gift_extract_grpid(code)
	if grpid == 0 {
		return nil
	}

	var info GiftInfo
	err := dbmgr.DBShare.GetProjectionByCond(
		dbmgr.C_tabname_giftinfo,
		db.M{
			"_id":   grpid,
			"codes": code,
		},
		db.M{"memo": 0, "codes": 0},
		&info,
	)
	if err != nil {
		return nil
	}

	return &info
}

func gift_code_used(info *GiftInfo, code string, uid string) bool {
	// (grpid, uid) || reuse ? (code, uid) : (code)
	cond := make([]db.M, 0, 2)

	cond = append(cond, db.M{"grpid": info.GrpId, "userid": uid})
	if info.Reuse != 0 {
		cond = append(cond, db.M{"code": code, "userid": uid})
	} else {
		cond = append(cond, db.M{"code": code})
	}

	var obj db.M
	err := dbmgr.DBShare.GetObjectByCond(
		dbmgr.C_tabname_giftuse,
		db.M{"$or": cond},
		&obj,
	)
	if err == nil {
		return true
	} else if db.IsNotFound(err) {
		return false
	} else {
		return true
	}
}

func gift_update_use(info *GiftInfo, code string, uid string) {
	grpid := info.GrpId

	async.Push(func() {
		err := dbmgr.DBShare.Insert(
			dbmgr.C_tabname_giftuse,
			db.M{
				"grpid":  grpid,
				"area":   info.Area,
				"code":   code,
				"userid": uid,
				"svr":    config.CurGame.Id,
				"ts":     time.Now(),
			},
		)
		if err != nil {
			log.Error("flush giftuse failed:", grpid, code, uid, err)
		}
	})
}
