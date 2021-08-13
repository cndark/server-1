package guild

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app/comp"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

const (
	C_guild_max_plr_applies = 5
	C_guild_max_applies     = 50
)

// 申请模式
const (
	AM_Accept = 1
	AM_Apply  = 2
	AM_Deny   = 3
)

// 离会原因
const (
	LR_Leave   = 1
	LR_Kick    = 2
	LR_Destroy = 3
)

// ============================================================================

type Guild struct {
	Id        string             `bson:"-"`
	Name      string             ``
	Icon      int32              ``
	Svr       string             ``
	Exp       int64              ``
	Lv        int32              ``
	CreateTs  time.Time          ``
	CreatePlr string             ``
	Members   map[string]*Member `bson:"mb"`
	Notice    string             `` // 公告
	Apply     *apply_t           ``
	Log       []*log_t           ``
	Wish      map[int64]*wish_t  // 祈愿 [seq]
	Harbor    *harbor_t          // 港口
	Boss      *boss_t            // 副本
	ZmTs      time.Time          // 发布招募的时间

	owner  *Member // guild owner
	atkpwr int32   // 战力
	rank   int32   // 排名

	// save
	save_tid *core.Timer

	db *db.Database
}

type apply_t struct {
	Mode   int32    ``
	NeedLv int32    ``
	AList  []string `bson:"alist,omitempty"`
}

// ============================================================================
// apply_t

func (self *apply_t) add(plrid string) {
	self.AList = append(self.AList, plrid)
	if len(self.AList) > C_guild_max_applies {
		self.AList = self.AList[1:]
	}
}

func (self *apply_t) remove(plrid string) {
	for i, v := range self.AList {
		if v == plrid {
			self.AList = append(self.AList[:i], self.AList[i+1:]...)
			return
		}
	}
}

func (self *apply_t) clone_list() (ret []string) {
	ret = make([]string, 0, len(self.AList))

	for _, v := range self.AList {
		ret = append(ret, v)
	}

	return
}

func (self *apply_t) ToMsg() (ret []*msg.GuildApplyRow) {
	ret = make([]*msg.GuildApplyRow, 0, len(self.AList))

	for _, plrid := range self.AList {
		plr := find_player(plrid)
		if plr != nil {
			ret = append(ret, &msg.GuildApplyRow{
				Plr: plr.ToMsg_SimpleInfo(),
			})
		}
	}

	return
}

// ============================================================================
// 存盘

func (self *Guild) open() {
	// bind
	for _, m := range self.Members {
		m.gld = self

		if m.Rank == RK_Owner {
			self.owner = m
		}
	}

	// init sub sys of guild
	self.Harbor.init(self)
	self.Boss.init(self)

	// start save
	self.save_timer_start()
}

func (self *Guild) close() {
	// stop save
	self.save_timer_stop()

	// -------- final save --------
	self.save()
}

func (self *Guild) destroy() bool {
	// remove from db
	err := self.DB().Remove(dbmgr.C_tabname_guild, self.Id)
	if err != nil {
		log.Error("destroy guild failed:", self.Id, err)
		return false
	}

	// stop save
	self.save_timer_stop()

	return true
}

func (self *Guild) save_timer_start() {
	self.save_tid = loop.SetTimeout(time.Now().Add(time.Duration(1200+rand.Intn(1200))*time.Second), func() {
		self.auto_kick_nonactive()
		self.save_async()
		self.save_timer_start()
	})
}

func (self *Guild) save_timer_stop() {
	if self.save_tid != nil {
		loop.CancelTimer(self.save_tid)
		self.save_tid = nil
	}
}

// 创建时存盘：仅用于创建
func (self *Guild) save_create() bool {
	err := self.DB().Insert(
		dbmgr.C_tabname_guild,
		db.M{"_id": self.Id, "base": self},
	)
	if err != nil {
		log.Error("guild save_create failed:", err)
		return false
	}

	return true
}

// 异步存盘：用于定时存盘
func (self *Guild) save_async() {
	// clone
	obj := core.CloneBsonObject(self)

	// async save
	async.Push(
		func() {
			self.DB().Update(
				dbmgr.C_tabname_guild,
				self.Id,
				db.M{"$set": db.M{"base": obj}},
			)
		},
	)
}

// 同步存盘：用于停服时存盘
func (self *Guild) save() {
	self.DB().Update(
		dbmgr.C_tabname_guild,
		self.Id,
		db.M{"$set": db.M{"base": self}},
	)
}

func (self *Guild) DB() *db.Database {
	return self.db
}

// ============================================================================

func (self *Guild) Owner() *Member {
	return self.owner
}

func (self *Guild) GetId() string {
	return self.Id
}

func (self *Guild) GetLevel() int32 {
	return self.Lv
}

func (self *Guild) GetName() string {
	return self.Name
}

func (self *Guild) GetAtkPwr() int32 {
	return self.atkpwr
}

func (self *Guild) GetRank() int32 {
	return self.rank
}

func (self *Guild) FindMember(plrid string) *Member {
	return self.Members[plrid]
}

func (self *Guild) AddExp(plr comp.IPlayer, v int64) {
	if v <= 0 {
		return
	}

	new_lv := self.Lv
	self.Exp += v

	// consume exp
	for {
		conf := gamedata.ConfGuild.Query(new_lv)
		if conf == nil || conf.UpgradeExp == 0 {
			// 满级
			self.Exp = 0
			break
		}

		if self.Exp < conf.UpgradeExp {
			break
		}

		self.Exp -= conf.UpgradeExp

		new_lv++
	}

	// fire
	evtmgr.Fire(gconst.Evt_GuildExpAdd, self, v)
	evtmgr.Fire(gconst.Evt_GuildExpPlrAdd, plr, v)

	// check levelup
	if new_lv == self.Lv {
		self.Broadcast(&msg.GS_Guild_Lv{
			Level: -1, // 无变化
			Exp:   self.Exp,
		})

		return
	}

	// set new lv
	old := self.Lv
	self.Lv = new_lv

	// event: levelup
	self.OnLevelup(old, new_lv)

	evtmgr.Fire(gconst.Evt_GuildChange, self, plr)
}

func (self *Guild) OnLevelup(old_lv, new_lv int32) {
	// notify
	self.Broadcast(&msg.GS_Guild_Lv{
		Level: new_lv,
		Exp:   self.Exp,
	})

	// fire
	evtmgr.Fire(gconst.Evt_GuildLv, self, new_lv)
}

func (self *Guild) AddMember(plr IPlayer) int32 {
	if self.IsFull() {
		return Err.Guild_FullMember
	}

	plrid := plr.GetId()

	if self.Members[plrid] != nil {
		return Err.Failed
	}

	// new member
	m := &Member{
		Id:   plrid,
		Rank: RK_Member,

		gld: self,
	}

	// add
	self.Members[plrid] = m

	// bind
	plr.BindGuild(self.Id)

	// purge apply-list
	GuildMgr.alist_purge_for_plr(plrid)

	// notify
	self.Broadcast(&msg.GS_Guild_Join{
		GuildId:   self.Id,
		GuildName: self.Name,
		Mb:        m.ToMsg_Info(),
	})

	// fire
	evtmgr.Fire(gconst.Evt_GuildJoin, plr, self)
	evtmgr.Fire(gconst.Evt_GuildUserChange, self, plr)

	// log
	self.AddLog(C_Guild_Log_Join, map[string]string{
		"player": plr.GetName(),
	})

	return Err.OK
}

func (self *Guild) RemoveMember(plr IPlayer, reason int32) {
	plrid := plr.GetId()

	m := self.Members[plrid]
	if m == nil {
		return
	}

	// unbind
	plr.BindGuild("")

	// notify
	self.Broadcast(&msg.GS_Guild_Leave{
		Reason: reason,
		PId:    plr.GetId(),
		PName:  plr.GetName(),
	})

	// remove
	delete(self.Members, plrid)

	// fire
	evtmgr.Fire(gconst.Evt_GuildLeave, self, plr, reason)
	evtmgr.Fire(gconst.Evt_GuildUserChange, self, plr)

	// log
	self.AddLog(C_Guild_Log_Leave, map[string]string{
		"player": plr.GetName(),
	})

	// email
	if reason == LR_Kick {
		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return
		}

		conf_m := gamedata.ConfMail.Query(conf_g.GuildExitMail)
		if conf_m == nil {
			return
		}

		ml := mail.New(plr).SetKey(conf_g.GuildExitMail)
		ml.AddDict("gldName", self.GetName())
		for _, k := range conf_m.MailItem {
			ml.AddAttachment(k.Id, float64(k.N))
		}

		ml.Send()
	}
}

func (self *Guild) IsFull() bool {
	conf := gamedata.ConfGuild.Query(self.Lv)
	if conf == nil {
		return true
	}

	return int32(len(self.Members)) >= conf.MbLimit
}

func (self *Guild) IsFullForRank(rk int32) int32 {
	conf := gamedata.ConfGuild.Query(self.Lv)
	if conf == nil {
		return Err.Guild_FullMember
	}

	n := int32(0)
	for _, m := range self.Members {
		if m.Rank == rk {
			n++
		}
	}

	switch rk {
	case RK_Vice:
		if n >= conf.ViceLimit {
			return Err.Guild_FullVice
		}
	default:
	}

	return Err.OK
}

func (self *Guild) ChangeName(name string, f func(bool)) {
	oldname := self.Name

	async.Push(
		func() {
			// update name-db
			if !dbmgr.Center_ChangeGName(oldname, name) {
				loop.Push(func() {
					f(false)
				})
				return
			}

			// update guild-db
			err := self.DB().Update(
				dbmgr.C_tabname_guild,
				self.Id,
				db.M{"$set": db.M{"base.name": name}},
			)
			if err != nil {
				log.Warning("Guild.ChangeName() failed:", err)
			}

			// update memory
			loop.Push(func() {
				self.Name = name
				delete(GuildMgr.gld_by_name, oldname)
				GuildMgr.gld_by_name[name] = self

				f(true)
			})
		},
	)
}

func (self *Guild) update_atkpwr() {
	self.atkpwr = 0

	for _, m := range self.Members {
		plr := load_player(m.Id)
		if plr != nil {
			self.atkpwr += plr.GetAtkPwr()
		}
	}
}

func (self *Guild) daily_reset() {
	for _, m := range self.Members {
		m.daily_reset()
	}
}

// ============================================================================
// msg
func (self *Guild) Broadcast(message msg.Message, except ...string) {
	var exceptid string

	if len(except) > 0 {
		exceptid = except[0]
	}

	for _, m := range self.Members {
		if m.Id == exceptid {
			continue
		}

		plr := find_player(m.Id)
		if plr != nil {
			plr.SendMsg(message)
		}
	}
}

func (self *Guild) broadcast_rank(rk int32, message msg.Message) {
	for _, m := range self.Members {
		if m.Rank > rk {
			continue
		}

		plr := find_player(m.Id)
		if plr != nil {
			plr.SendMsg(message)
		}
	}
}

func (self *Guild) ToMsg_Row() *msg.GuildRow {
	return &msg.GuildRow{
		Id:      self.Id,
		Name:    self.Name,
		Icon:    self.Icon,
		Lv:      self.Lv,
		MemberN: int32(len(self.Members)),
		NeedLv:  self.Apply.NeedLv,
		Rank:    self.rank,
		AtkPwr:  self.atkpwr,
	}
}

func (self *Guild) ToMsg_InfoFull(plr IPlayer) *msg.GuildInfo_Full {
	info := &msg.GuildInfo_Full{
		Row:         self.ToMsg_Row(),
		Exp:         self.Exp,
		Notice:      self.Notice,
		AMode:       self.Apply.Mode,
		HarborLevel: self.Harbor.Lv,
		HarborXp:    self.Harbor.Xp,
	}

	info.Mbs = make([]*msg.GuildMemberInfo, 0, len(self.Members))
	for _, m := range self.Members {
		info.Mbs = append(info.Mbs, m.ToMsg_Info())
	}

	if self.owner.Id == plr.GetId() {
		info.ZmTs = self.ZmTs.Unix()
	}

	info.PlrData = plr.GetGuildPlrData().ToMsg()

	return info
}

func (self *Guild) ToMsg_Log() []*msg.GuildLog {
	ret := []*msg.GuildLog{}
	for _, v := range self.Log {
		ret = append(ret, &msg.GuildLog{
			Id:    v.Id,
			Param: v.Param,
			Ts:    v.Ts.Unix(),
		})
	}

	return ret
}

// ============================================================================
// 其他

func (self *Guild) LockMemberLeave() bool {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return false
	}

	if len(conf.NotExitGuild) != 2 {
		return false
	}

	now := time.Now()
	t1, t2 := core.ParseTime(conf.NotExitGuild[0]), core.ParseTime(conf.NotExitGuild[1])

	if now.After(t1) && now.Before(t2) {
		return true
	}

	return false
}

func (self *Guild) auto_kick_nonactive() {
	// find
	var arr []*Member

	for _, m := range self.Members {
		plr := load_player(m.Id)
		if plr == nil {
			continue
		}

		if !plr.IsActive() && m.Rank != RK_Owner {
			arr = append(arr, m)
		}
	}

	// kick
	for _, m := range arr {
		m.Leave(LR_Leave)
	}
}

// ============================================================================
// 入会

func (self *Guild) SetApplyMode(mode int32, need_lv int32) int32 {
	// set need lv
	if need_lv >= 1 {
		self.Apply.NeedLv = need_lv
	}

	// set mode
	if mode >= AM_Accept && mode <= AM_Deny {
		self.Apply.Mode = mode
	}

	// check
	switch self.Apply.Mode {
	case AM_Accept:
		self.ApplyAcceptOnekey()

	case AM_Apply:
		// do nothing

	case AM_Deny:
		self.ApplyDenyOneKey()
	}

	return Err.OK
}

func (self *Guild) SetIcon(v int32) int32 {
	self.Icon = v
	self.Broadcast(&msg.GS_Guild_Icon{
		Icon: v,
	})

	return Err.OK
}

func (self *Guild) SetNotice(v string) int32 {
	self.Notice = v
	self.Broadcast(&msg.GS_Guild_Notice{
		Notice: v,
	})

	return Err.OK
}

func (self *Guild) ApplyRequst(plr IPlayer) int32 {
	// plr already in guild ?
	if plr.GetGuild() != nil {
		return Err.Guild_PlrInGuild
	}

	// check mode
	switch self.Apply.Mode {
	case AM_Accept:
		if plr.GetLevel() < self.Apply.NeedLv {
			return Err.Guild_PlrLowLevel
		}

		return self.AddMember(plr)

	case AM_Apply:
		plrid := plr.GetId()

		// applied ?
		if GuildMgr.plr_alist_exists(plrid, self.Id) {
			return Err.Guild_AlreadyApplied
		}

		// add apply
		if !GuildMgr.plr_alist_add(plrid, self.Id) {
			return Err.Guild_FullPlrApply
		}
		self.Apply.add(plrid)

		// notify
		self.broadcast_rank(RK_Vice, &msg.GS_Guild_NewApply{
			PId: plrid,
		})

		return Err.OK

	case AM_Deny:
		return Err.Guild_JoinDenied

	default:
		return Err.Guild_InvalidApplyMode
	}
}

func (self *Guild) ApplyCancel(plr IPlayer) int32 {
	return self.ApplyDeny(plr.GetId())
}

func (self *Guild) ApplyAccept(plrid string) int32 {
	// applied ?
	if !GuildMgr.plr_alist_exists(plrid, self.Id) {
		return Err.Guild_NotApplied
	}

	plr := load_player(plrid)
	if plr == nil {
		return Err.Failed
	}

	return self.AddMember(plr)
}

func (self *Guild) ApplyDeny(plrid string) int32 {
	// applied ?
	if !GuildMgr.plr_alist_exists(plrid, self.Id) {
		return Err.Guild_NotApplied
	}

	// remove apply
	GuildMgr.plr_alist_remove(plrid, self.Id)
	self.Apply.remove(plrid)

	return Err.OK
}

func (self *Guild) ApplyAcceptOnekey() {
	// audit apply list before purge
	for _, plrid := range self.Apply.clone_list() {
		plr := find_player(plrid)
		if plr != nil && plr.GetLevel() >= self.Apply.NeedLv {
			if self.AddMember(plr) != Err.OK {
				break
			}
		}
	}

	// purge guild apply-list
	GuildMgr.alist_purge_for_guild(self)
}

func (self *Guild) ApplyDenyOneKey() {
	// remove all apply
	GuildMgr.alist_purge_for_guild(self)
}
