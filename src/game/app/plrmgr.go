package app

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// ============================================================================

var PlayerMgr = &plrmgr_t{
	plrs_by_id:   make(map[string]*Player),
	plrs_by_name: make(map[string]*Player),
	plrs_online:  make(map[uint64]*Player),
}

// ============================================================================

type plrmgr_t struct {
	plrs_by_id   map[string]*Player // loaded players by id
	plrs_by_name map[string]*Player // loaded players by name
	plrs_online  map[uint64]*Player // online players by sid

	num_loaded int32 // loaded player number
	num_online int32 // online player number
	num_reg    int32 // register player number
}

// ============================================================================

func get_user_dbname(uid string) string {
	return strings.Split(uid, "-")[0]
}

// ============================================================================

func (self *plrmgr_t) Open() {
	self.load_active_players()
	self.load_num_reg()
}

func (self *plrmgr_t) Close() {
	c := core.NewConsumer(8)

	L := len(self.plrs_by_id)
	i := 0
	j := 0
	a := []int{0, 10, 33, 50, 66, 80, 100}

	for _, plr := range self.plrs_by_id {
		c.Push(plr.save)

		i++
		pct := i * 100 / L
		if pct >= a[j] {
			j++
			log.Infof("saving players progress: %d%%", pct)
		}
	}

	c.Close()

	log.Info("ALL players are saved")
}

func (self *plrmgr_t) LoadPlayer(uid string) *Player {
	plr, _ := self.load_player(uid)
	return plr
}

func (self *plrmgr_t) LoadPlayerWithError(uid string) (*Player, int32) {
	return self.load_player(uid)
}

func (self *plrmgr_t) LoadIPlayer(uid string) interface{} {
	plr, _ := self.load_player(uid)
	if plr == nil {
		return nil
	} else {
		return plr
	}
}

func (self *plrmgr_t) CreatePlayer(uid string, f func(*User)) *Player {
	user := create_user(uid, f)
	if user == nil {
		return nil
	}

	// new player
	plr := new_player(user)

	// add to mgr
	if !self.add_player(plr, true) {
		plr = nil
	}

	// update num reg
	self.num_reg++

	// fire
	evtmgr.Fire(gconst.Evt_ServerNewUser, self.num_reg, plr)

	return plr
}

func (self *plrmgr_t) SetOnline(plr *Player, sid uint64) {
	if plr.sid != 0 || sid == 0 {
		return
	}

	// set sid
	plr.sid = sid

	// add to mgr
	self.plrs_online[sid] = plr

	// event: online
	plr.OnOnline()

	// count
	atomic.AddInt32(&self.num_online, 1)
}

func (self *plrmgr_t) SetOffline(plr *Player, shutdown bool) {
	if plr.sid == 0 {
		return
	}

	// event: offline
	plr.OnOffline(shutdown)

	// remove from mgr
	delete(self.plrs_online, plr.sid)

	// reset sid
	plr.sid = 0

	// count
	atomic.AddInt32(&self.num_online, -1)
}

func (self *plrmgr_t) OfflinePlayers(gateid int32) {
	if gateid == 0 {
		// all players
		for _, plr := range self.plrs_by_id {
			self.SetOffline(plr, true)
		}
	} else {
		// players from gateid
		for _, plr := range self.plrs_by_id {
			if sid_to_gateid(plr.sid) == gateid {
				self.SetOffline(plr, true)
			}
		}
	}
}

func (self *plrmgr_t) FindPlayerById(uid string) *Player {
	return self.plrs_by_id[uid]
}

func (self *plrmgr_t) FindIPlayerById(uid string) interface{} {
	plr := self.plrs_by_id[uid]
	if plr == nil {
		return nil
	} else {
		return plr
	}
}

func (self *plrmgr_t) FindPlayerByName(name string) *Player {
	return self.plrs_by_name[name]
}

func (self *plrmgr_t) FindIPlayerByName(name string) interface{} {
	plr := self.plrs_by_name[name]
	if plr == nil {
		return nil
	} else {
		return plr
	}
}

func (self *plrmgr_t) FindPlayerBySid(sid uint64) *Player {
	return self.plrs_online[sid]
}

func (self *plrmgr_t) ForEachLoadedIPlayer(f func(plr interface{})) {
	for _, plr := range self.plrs_by_id {
		f(plr)
	}
}

func (self *plrmgr_t) ForEachOnlineIPlayer(f func(plr interface{})) {
	for _, plr := range self.plrs_online {
		f(plr)
	}
}

func (self *plrmgr_t) ForEachLoadedIPlayerBreakable(f func(plr interface{}) bool) {
	for _, plr := range self.plrs_by_id {
		if !f(plr) {
			break
		}
	}
}

func (self *plrmgr_t) ForEachOnlineIPlayerBreakable(f func(plr interface{}) bool) {
	for _, plr := range self.plrs_online {
		if !f(plr) {
			break
		}
	}
}

func (self *plrmgr_t) UpdatePlayerName(plr *Player, name string) {
	if name == plr.user.Name {
		return
	}

	delete(self.plrs_by_name, plr.user.Name)
	plr.user.Name = name
	self.plrs_by_name[name] = plr
}

func (self *plrmgr_t) NumLoaded() int32 {
	return atomic.LoadInt32(&self.num_loaded)
}

func (self *plrmgr_t) NumOnline() int32 {
	return atomic.LoadInt32(&self.num_online)
}

// ============================================================================

func (self *plrmgr_t) load_active_players() {
	log.Info("loading active players ...")

	now := time.Now()
	ts := now.Add(-gconst.PLAYER_ActiveDays * 24 * time.Hour)

	for k, _ := range config.Common.DBUser {
		udb := dbmgr.UserDB(k)
		if udb == nil {
			core.Panic("udb NOT found:", k)
		}

		var arr []*struct {
			Id   string `bson:"_id"`
			Base *User
		}

		err := udb.GetAllProjectionsByCond(
			dbmgr.C_tabname_user,
			db.M{
				"base.svr":      config.CurGame.Name,
				"base.login_ts": db.M{"$gt": ts},
				"base.ban_ts":   db.M{"$lt": now},
			},
			db.M{"base": 1},
			&arr,
		)
		if err != nil {
			core.Panic("loading players failed:", err)
		}

		for _, obj := range arr {
			// bind
			obj.Base.Id = obj.Id
			obj.Base.db = udb

			// add
			plr := new_player(obj.Base)
			self.add_player(plr, false)
		}
	}

	runtime.GC()

	log.Infof("    %d players are loaded", len(self.plrs_by_id))
}

func (self *plrmgr_t) load_num_reg() {
	n, err := dbmgr.DBShare.Count(
		dbmgr.C_tabname_userinfo,
		db.M{"area": config.Common.Area.Id, "svr": config.CurGame.Name},
	)
	if err != nil {
		core.Panic("loading register player count failed:", err)
	}

	self.num_reg = int32(n)
}

func (self *plrmgr_t) load_player(uid string) (plr *Player, ec int32) {
	//check
	if uid == "" {
		return nil, Err.Failed
	}

	// find in memory
	plr = self.plrs_by_id[uid]
	if plr != nil {
		return plr, Err.OK
	}

	// load user from db
	user, ec := self.load_from_db(uid)
	if ec == Err.OK {
		// new player
		plr = new_player(user)

		// add to mgr
		if self.add_player(plr, false) {
			return plr, Err.OK
		} else {
			return nil, Err.Failed
		}
	}

	// failed
	return nil, ec
}

func (self *plrmgr_t) load_from_db(uid string) (*User, int32) {
	// get user db
	dbname := get_user_dbname(uid)
	udb := dbmgr.UserDB(dbname)
	if udb == nil {
		log.Critical("get user db failed:", dbname)
		log.Critical(core.Callstack())
		return nil, Err.Failed
	}

	// load
	var obj struct {
		Base *User
	}
	err := udb.GetProjectionByCond(
		dbmgr.C_tabname_user,
		db.M{
			"_id":      uid,
			"base.svr": config.CurGame.Name,
		},
		db.M{"base": 1},
		&obj,
	)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, Err.Plr_NotFound
		} else {
			log.Critical("load user from db failed:", err)
			log.Critical(core.Callstack())
			return nil, Err.Failed
		}
	}

	// bind
	obj.Base.Id = uid
	obj.Base.db = udb

	// return
	return obj.Base, Err.OK
}

func (self *plrmgr_t) add_player(plr *Player, creation bool) bool {
	// open: if data-loading fails at open phase, we'd rather report error instead of going on
	if !plr.open() {
		return false
	}

	self.plrs_by_id[plr.user.Id] = plr
	self.plrs_by_name[plr.user.Name] = plr

	// loaded
	plr.loaded()

	// created
	if creation {
		plr.created()
	}

	// patch
	plr.patch()

	// count
	atomic.AddInt32(&self.num_loaded, 1)

	return true
}
