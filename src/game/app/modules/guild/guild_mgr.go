package guild

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"sort"
	"time"
)

// ============================================================================

var GuildMgr = &guild_mgr_t{
	gld_by_id:   make(map[string]*Guild),
	gld_by_name: make(map[string]*Guild),
	gld_arr:     nil,

	plr_alist: make(map[string]map[string]bool),
}

// ============================================================================

const (
	C_guild_list_rpp = 10
)

// ============================================================================

type guild_mgr_t struct {
	gld_by_id   map[string]*Guild // by id
	gld_by_name map[string]*Guild // by name
	gld_arr     []*Guild          // sorted by total sp

	plr_alist map[string]map[string]bool // [plrid][gldid]
}

// ============================================================================
func init() {
	evtmgr.On(gconst.Evt_GlobalResetDaily, func(args ...interface{}) {
		GuildMgr.daily_reset()
	})
}

// ============================================================================

func (self *guild_mgr_t) Open() {
	// load guilds
	self.load_guilds()

	// start sort timer
	self.sort_timer_start()
}

func (self *guild_mgr_t) Close() {
	log.Info("saving guilds ...")

	for _, gld := range self.gld_by_id {
		gld.close()
	}

	log.Info("ALL guilds are saved")
}

func (self *guild_mgr_t) load_guilds() {
	log.Info("loading guilds ...")

	// load
	for k := range config.Common.DBUser {
		udb := dbmgr.UserDB(k)
		if udb == nil {
			continue
		}

		var objs []*struct {
			Id   string `bson:"_id"`
			Base *Guild
		}

		err := udb.GetAllProjectionsByCond(
			dbmgr.C_tabname_guild,
			db.M{"base.svr": config.CurGame.Name},
			db.M{"base": 1},
			&objs,
		)
		if err != nil {
			log.Warning("loading guild failed:", k, err)
			continue
		}

		for _, obj := range objs {
			// bind
			obj.Base.Id = obj.Id
			obj.Base.db = udb

			self.add(obj.Base, false)
			self.plr_alist_build(obj.Base)
		}
	}

	// delay the sort to avoid player loading.
	//	!Note: guilds are loaded before players
	loop.Push(func() {
		self.sort()
	})

	log.Infof("    %d guilds are loaded", len(self.gld_by_id))
}

func (self *guild_mgr_t) CreateGuild(plr IPlayer, name string, notice string, icon int32) (ec int32, gld *Guild) {
	// alloc db
	dbname := dbmgr.Center_AllocUserDB()
	if dbname == "" {
		ec = Err.Failed
		return
	}
	udb := dbmgr.UserDB(dbname)
	if udb == nil {
		log.Error("get user db failed:", dbname)
		log.Error(core.Callstack())
		ec = Err.Failed
		return
	}

	// register guild name
	if !dbmgr.Center_InsertGName(name) {
		ec = Err.Guild_DupName
		return
	}

	// update udb load
	dbmgr.Center_IncUserLoad(dbname)

	// new guild
	gld = &Guild{
		Id:        dbmgr.Center_GenGuildId(dbname),
		Name:      name,
		Icon:      icon,
		Svr:       config.CurGame.Name,
		Exp:       0,
		Lv:        1,
		CreateTs:  time.Now(),
		CreatePlr: plr.GetId(),
		Members:   make(map[string]*Member),
		Notice:    notice,
		Apply: &apply_t{
			Mode:   AM_Accept,
			NeedLv: 1,
			AList:  nil,
		},
		Log:    nil,
		Wish:   make(map[int64]*wish_t),
		Harbor: new_harbor(),
		Boss:   new_boss(),

		db: udb,
	}

	// save on create
	if !gld.save_create() {
		ec = Err.Failed
		return
	}

	// add
	self.add(gld, true)

	// add owner
	gld.AddMember(plr)
	gld.FindMember(plr.GetId()).SetRank(RK_Owner)

	// fire
	evtmgr.Fire(gconst.Evt_GuildCreate, gld, plr)

	gld.AddLog(C_Guild_Log_Create, map[string]string{
		"player":  plr.GetName(),
		"gldName": gld.Name,
	})

	ec = Err.OK
	return
}

func (self *guild_mgr_t) DestroyGuild(gld *Guild) int32 {
	// destroy
	if !gld.destroy() {
		return Err.Failed
	}

	// fire
	evtmgr.Fire(gconst.Evt_GuildDestroy, gld, load_player(gld.owner.Id))

	// dismiss members
	conf := gamedata.ConfGlobalPublic.Query(1)

	arr := make([]*Member, 0, len(gld.Members))
	for _, m := range gld.Members {
		arr = append(arr, m)
	}

	for _, m := range arr {
		plr := load_player(m.Id)
		if plr == nil {
			continue
		}

		gld.RemoveMember(plr, LR_Destroy)

		if conf != nil {
			self.destroy_mail(plr, conf.GuildDisMail)
		}
	}

	// remove
	self.remove(gld)

	return Err.OK
}

func (self *guild_mgr_t) FindGuild(id string) *Guild {
	return self.gld_by_id[id]
}

func (self *guild_mgr_t) FindGuildByName(name string) *Guild {
	return self.gld_by_name[name]
}

func (self *guild_mgr_t) Array_Guilds() []*Guild {
	return self.gld_arr
}

func (self *guild_mgr_t) FetchGuildList(page int32) (ret []*msg.GuildRow) {
	// check page
	if page < 1 {
		page = 1
	}

	// calc paging
	L := int32(len(self.gld_arr))

	a := (page - 1) * C_guild_list_rpp
	b := a + C_guild_list_rpp

	if a >= L {
		return
	}
	if b > L {
		b = L
	}

	// fetch
	ret = make([]*msg.GuildRow, 0, C_guild_list_rpp)
	for _, gld := range self.gld_arr[a:b] {
		ret = append(ret, gld.ToMsg_Row())
	}

	return
}

func (self *guild_mgr_t) FetchGuildList_PlayerApplied(plrid string) (ret []*msg.GuildRow) {
	// fetch
	for gldid := range self.plr_alist[plrid] {
		gld := self.FindGuild(gldid)
		if gld != nil {
			ret = append(ret, gld.ToMsg_Row())
		}
	}

	return
}

// ============================================================================
func (self *guild_mgr_t) daily_reset() {
	for _, gld := range self.gld_by_id {
		gld.daily_reset()
	}
}

func (self *guild_mgr_t) add(gld *Guild, srt bool) {
	if self.gld_by_id[gld.Id] != nil {
		return
	}

	self.gld_by_id[gld.Id] = gld
	self.gld_by_name[gld.Name] = gld
	self.gld_arr = append(self.gld_arr, gld)

	// open
	gld.open()

	// sort
	if srt {
		self.sort()
	}
}

func (self *guild_mgr_t) remove(gld *Guild) {
	// remove
	delete(self.gld_by_id, gld.Id)
	delete(self.gld_by_name, gld.Name)

	for i, v := range self.gld_arr {
		if v.Id == gld.Id {
			L := len(self.gld_arr)
			self.gld_arr[i] = self.gld_arr[L-1]
			self.gld_arr = self.gld_arr[:L-1]
			break
		}
	}

	// sort
	self.sort()
}

func (self *guild_mgr_t) sort_timer_start() {
	// total speed sort
	{
		var f func()
		f = func() {
			loop.SetTimeout(time.Now().Add(5*time.Minute), func() {
				self.sort()
				f()
			})
		}
		f()
	}
}

func (self *guild_mgr_t) sort() {
	// update guild speedr
	for _, gld := range self.gld_by_id {
		gld.update_atkpwr()
	}

	// sort
	sort.Slice(self.gld_arr, func(i, j int) bool {
		return self.gld_arr[i].GetAtkPwr() > self.gld_arr[j].GetAtkPwr()
	})

	// update rank
	for i, gld := range self.gld_arr {
		gld.rank = int32(i) + 1
	}
}

// ============================================================================
// 申请列表

func (self *guild_mgr_t) plr_alist_build(gld *Guild) {
	for _, plrid := range gld.Apply.AList {
		self.plr_alist_add(plrid, gld.Id)
	}
}

func (self *guild_mgr_t) plr_alist_add(plrid string, gldid string) bool {
	lst := self.plr_alist[plrid]
	if lst == nil {
		lst = make(map[string]bool)
		self.plr_alist[plrid] = lst
	}

	if len(lst) >= C_guild_max_plr_applies {
		return false
	}

	lst[gldid] = true

	return true
}

func (self *guild_mgr_t) plr_alist_remove(plrid string, gldid string) {
	lst := self.plr_alist[plrid]
	if lst != nil {
		delete(lst, gldid)
	}
}

func (self *guild_mgr_t) plr_alist_exists(plrid string, gldid string) bool {
	return self.plr_alist[plrid][gldid]
}

// 销毁玩家所有申请 (双向)
func (self *guild_mgr_t) alist_purge_for_plr(plrid string) {
	for gldid := range self.plr_alist[plrid] {
		gld := self.FindGuild(gldid)
		if gld != nil {
			gld.Apply.remove(plrid)
		}
	}

	delete(self.plr_alist, plrid)
}

// 销毁公会所有申请 (双向)
func (self *guild_mgr_t) alist_purge_for_guild(gld *Guild) {
	for _, plrid := range gld.Apply.AList {
		self.plr_alist_remove(plrid, gld.Id)
	}

	gld.Apply.AList = nil
}

func (self *guild_mgr_t) destroy_mail(plr IPlayer, mid int32) {
	conf_m := gamedata.ConfMail.Query(mid)
	if conf_m == nil {
		return
	}

	// new mail
	m := mail.New(plr).SetKey(mid)

	for _, k := range conf_m.MailItem {
		m.AddAttachment(k.Id, float64(k.N))
	}

	// send
	m.Send()
}
