package evt

import (
	"fw/src/core/sched/async"
	"fw/src/game/app"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

type rec_t struct {
	tp string
	d  map[string]interface{}
}

// ============================================================================

func new_benu_rec(op string, plr *app.Player) *rec_t {
	rec := &rec_t{
		tp: "glog_benu",
		d: map[string]interface{}{
			"op":    op,
			"area":  config.Common.Area.Id,
			"svrid": config.CurGame.Id,
			"ts":    time.Now(),
		},
	}

	if plr != nil {
		rec.d["authid"] = plr.GetAuthId()
		rec.d["uid"] = plr.GetId()
		rec.d["sdk"] = plr.GetSdk()
		rec.d["cts"] = plr.GetCreateTs()
		rec.d["vip"] = plr.GetVipLevel()
	}

	return rec
}

// func new_chuxin_rec(op string, fc string, plr *app.Player) *rec_t {
// 	rec := &rec_t{
// 		tp: "glog_chuxin",
// 		d: map[string]interface{}{
// 			"op":          op,
// 			"func":        fc,
// 			"platform_id": config.Common.Area.Id,
// 			"server_id":   config.CurGame.Id,
// 			"ts":          time.Now(),
// 		},
// 	}

// 	if plr != nil {
// 		rec.d["account"] = plr.GetAuthId()
// 		rec.d["role_id"] = plr.GetId()
// 		rec.d["channel_id"] = plr.GetSdk()
// 		rec.d["cts"] = plr.GetCreateTs()
// 	}

// 	return rec
// }

func (self *rec_t) set(k string, v interface{}) *rec_t {
	self.d[k] = v
	return self
}

func (self *rec_t) insert() {
	if dbmgr.Kfk != nil {
		// insert into kafka
		async.PushQ(gconst.AQ_GLog, func() {
			dbmgr.Kfk.Producer(self.tp, self.d)
		})
	} else {
		// insert into mongodb
		async.PushQ(gconst.AQ_GLog, func() {
			dbmgr.DBLog.Insert(dbmgr.C_tabname_log, self.d)
		})
	}
}
