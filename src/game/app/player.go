package app

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/core/sched/next_tick"
	"fw/src/core/sched/resetter"
	"fw/src/game/app/comp"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/appoint"
	"fw/src/game/app/modules/arena"
	"fw/src/game/app/modules/attaintab"
	"fw/src/game/app/modules/bill"
	"fw/src/game/app/modules/billfirst"
	"fw/src/game/app/modules/chat"
	"fw/src/game/app/modules/cloud"
	"fw/src/game/app/modules/crusade"
	"fw/src/game/app/modules/daysign"
	"fw/src/game/app/modules/draw"
	"fw/src/game/app/modules/friend"
	"fw/src/game/app/modules/giftshop"
	"fw/src/game/app/modules/growfund"
	"fw/src/game/app/modules/guild"
	"fw/src/game/app/modules/gwar"
	"fw/src/game/app/modules/heroskin"
	"fw/src/game/app/modules/hframestore"
	"fw/src/game/app/modules/invite"
	"fw/src/game/app/modules/ladder"
	"fw/src/game/app/modules/lamp"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/marvelroll"
	"fw/src/game/app/modules/misc"
	"fw/src/game/app/modules/monthticket"
	"fw/src/game/app/modules/mopen"
	"fw/src/game/app/modules/privcard"
	"fw/src/game/app/modules/pushgift"
	"fw/src/game/app/modules/rank"
	"fw/src/game/app/modules/rift"
	"fw/src/game/app/modules/shop"
	"fw/src/game/app/modules/signdaily"
	"fw/src/game/app/modules/targetdays"
	"fw/src/game/app/modules/taskachv"
	"fw/src/game/app/modules/taskdaily"
	"fw/src/game/app/modules/taskgrow"
	"fw/src/game/app/modules/taskmonth"
	"fw/src/game/app/modules/teammgr"
	"fw/src/game/app/modules/tower"
	"fw/src/game/app/modules/tutorial"
	"fw/src/game/app/modules/vip"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/app/modules/wboss"
	"fw/src/game/app/modules/wlevel"
	"fw/src/game/app/modules/wleveldraw"
	"fw/src/game/app/modules/wlevelfund"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// ============================================================================

type Player struct {
	user *User
	sid  uint64

	hb_ts     time.Time // 上次心跳时间
	online_ts time.Time // 上次在线时长结算时间
}

// ============================================================================

func init() {
	// update player atk-power by teammgr
	evtmgr.On(gconst.Evt_DfdTeamUpdate, func(args ...interface{}) {
		plr := args[0].(*Player)

		plr.update_atkpwr()
	})
}

// ============================================================================
// 核心功能

func new_player(u *User) *Player {
	return &Player{
		user: u,
	}
}

func (self *Player) open() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("player open failed:", err)
		}
	}()

	// init
	user := self.user

	user.Bag.Init(self)
	user.Counter.Init(self) // MUST be called before other counter-related modules
	user.Bill.Init(self)
	user.AttainTab.Init(self)
	user.MailBox.Init(self)
	user.GuildPlrData.Init(self)
	user.MOpen.Init(self)
	user.Tutorial.Init(self)
	user.Cloud.Init(self)
	user.Misc.Init(self)
	user.Act.Init(self)
	user.TaskDaily.Init(self)
	user.TaskAchv.Init(self)
	user.Draw.Init(self)
	user.WLevel.Init(self)
	user.Appoint.Init(self)
	user.Tower.Init(self)
	user.TeamMgr.Init(self)
	user.Arena.Init(self)
	user.HFrameStore.Init(self)
	user.Shop.Init(self)
	user.MarvelRoll.Init(self)
	user.Friend.Init(self)
	user.Crusade.Init(self)
	user.Vip.Init(self)
	user.MonthTicket.Init(self)
	user.PushGift.Init(self)
	user.GiftShop.Init(self)
	user.PrivCard.Init(self)
	user.SignDaily.Init(self)
	user.TaskMonth.Init(self)
	user.DaySign.Init(self)
	user.RankPlay.Init(self)
	user.TargetDays.Init(self)
	user.WLevelFund.Init(self)
	user.GrowFund.Init(self)
	user.TaskGrow.Init(self)
	user.BillFirst.Init(self)
	user.GWar.Init(self)
	user.Rift.Init(self)

	if user.HeroSkin == nil {
		user.HeroSkin = heroskin.NewHeroSkin()
	}
	user.HeroSkin.Init(self)

	// 天梯
	if user.Ladder == nil {
		user.Ladder = ladder.NewLadder()
	}
	user.Ladder.Init(self)

	// 推图十连
	if user.WLevelDraw == nil {
		user.WLevelDraw = wleveldraw.NewWLevelDraw()
	}
	user.WLevelDraw.Init(self)

	// 本服杯赛
	if user.WarCup == nil {
		user.WarCup = warcup.NewWarCup()
	}
	user.WarCup.Init(self)

	// wboss
	if user.WBoss == nil {
		user.WBoss = wboss.NewWBoss()
	}
	user.WBoss.Init(self)

	// invite
	if user.Invite == nil {
		user.Invite = invite.NewInvite()
	}
	user.Invite.Init(self)

	// return
	return true
}

func (self *Player) loaded() {
	// bind guild
	self.BindGuild(self.GetGuildId())

	// hero apply mod props
	self.GetBag().ApplyHeroMods()

	// add to resetter
	resetter.Add(self)

	// atkpower
	self.update_atkpwr()

	// start save timer
	self.save_timer_start()
}

func (self *Player) created() {
	user := self.user

	// 设置 初始配置
	conf := gamedata.ConfInitial.Query(1)
	if conf != nil {
		// 初始玩家等级，经验
		user.Lv = conf.InitialPlayerLevel

		// new bagop
		op := user.Bag.NewOp(gconst.ObjFrom_Init)

		// 初始货币
		for _, v := range conf.InitialCurrency {
			op.Inc(v.Id, v.N)
		}

		// 初始道具
		for _, v := range conf.InitialItem {
			op.Inc(v.Id, int64(v.N))
		}

		// 初始英雄
		for _, v := range conf.InitialHero {
			op.Inc(v, 1)
		}

		// apply bagop
		op.Apply()

		// mail
		conf_m := gamedata.ConfMail.Query(conf.InitialMailId)
		if conf_m != nil {
			m := mail.New(self).SetKey(conf.InitialMailId)
			for _, v := range conf_m.MailItem {
				m.AddAttachment(v.Id, float64(v.N))
			}
			m.Send()
		}
	}

	// stats (initial bill-stats record)
	evtmgr.Fire(gconst.Evt_BillStats, self, int32(0), "", "", int32(0), "")

	// for bot: pressure test
	if config.Common.DevMode && self.user.Sdk == "soda.pressure" {
		self.bot_ptest_add_res()
	}

	// for bot: ai
	if self.user.Sdk == "soda.ai" {
		self.bot_ai_add_res()
	}
}

func (self *Player) bot_ptest_add_res() {
	// res
	op := self.user.Bag.NewOp(gconst.ObjFrom_Init)
	for _, v := range gamedata.ConfPressure.Items() {
		op.Inc(v.Id, v.N)
	}
	op.Apply()

	// level
	self.SetLevel(100)

	// vip
	self.GetVip().Lv = 5

	// wlevel
	self.GetWLevel().GMSetLevel(50)
}

func (self *Player) bot_ai_add_res() {
	op := self.user.Bag.NewOp(gconst.ObjFrom_Init)
	for _, v := range gamedata.ConfAiBot.Items() {
		if v.Type == 1 {
			// res
			op.Inc(v.Id, v.N)

		}
	}
	op.Apply()
}

func (self *Player) send_user_info() {
	user := self.user

	self.SendMsg(&msg.GS_UserInfo{
		UserId:    user.Id,
		Name:      user.Name,
		Head:      user.Head,
		HFrame:    user.HFrame,
		Lv:        user.Lv,
		Exp:       user.Exp,
		CreateTs:  user.CreateTs.Unix(),
		IsNew:     user.IsNew,
		LoginIP:   user.LoginIP,
		OnlineDur: user.OnlineDur,
		SvrTs:     time.Now().Unix(),

		GuildId: self.GetGuildId(),

		Invite: self.GetInvite().ToMsg(),

		ForLoading: &msg.ForLoading{},
		AtkPwr:     self.user.atkpwr,

		Bag:         self.GetBag().ToMsg(),
		Counter:     self.GetCounter().ToMsg(),
		AttainTab:   self.GetAttainTab().ToMsg(),
		Mails:       self.GetMailBox().ToMsg(),
		MOpen:       self.GetMOpen().ToMsg(),
		Tutorial:    self.GetTutorial().ToMsg(),
		Misc:        self.GetMisc().ToMsg(),
		TaskDaily:   self.GetTaskDaily().ToMsg(),
		TaskAchv:    self.GetTaskAchv().ToMsg(),
		Draw:        self.GetDraw().ToMsg(),
		WLevel:      self.GetWLevel().ToMsg(),
		Appoint:     self.GetAppoint().ToMsg(),
		Tower:       self.GetTower().ToMsg(),
		TeamMgr:     self.GetTeamMgr().ToMsg(),
		Arena:       self.GetArena().ToMsg(),
		HFrameStore: self.GetHFrameStore().ToMsg(),
		MarvelRoll:  self.GetMarvelRoll().ToMsg(),
		Friend:      self.GetFriend().ToMsg(),
		Vip:         self.GetVip().ToMsg(),
		Chat:        chat.ChatHistory_ToMsg(self),
		MonthTicket: self.GetMonthTicket().ToMsg(),
		PushGift:    self.GetPushGift().ToMsg(),
		GiftShop:    self.GetGiftShop().ToMsg(),
		PrivCard:    self.GetPrivCard().ToMsg(),
		SignDaily:   self.GetSignDaily().ToMsg(),
		TaskMonth:   self.GetTaskMonth().ToMsg(),
		DaySign:     self.GetDaySign().ToMsg(),
		RankPlay:    self.GetRankPlay().ToMsg(),
		TargetDays:  self.GetTargetDays().ToMsg(),
		TaskGrow:    self.GetTaskGrow().ToMsg(),
		WLevelFund:  self.GetWLevelFund().ToMsg(),
		GrowFund:    self.GetGrowFund().ToMsg(),
		BillFirst:   self.GetBillFirst().ToMsg(),
		Lamp:        lamp.Lamp_ToMsg(),
		Gwar:        self.GetGWar().ToMsg(),
		Rift:        rift.Rift_ToMsg(self),
		HeroSkin:    self.GetHeroSkin().ToMsg(),
		Ladder:      self.GetLadder().ToMsg(),
		WLevelDraw:  self.GetWLevelDraw().ToMsg(),
		WarCup:      self.GetWarCup().ToMsg(),
		WBoss:       self.GetWBoss().ToMsg(),
	})
}

func (self *Player) OnOnline() {
	log.Info(self.user.Id, self.user.Name, "is online")

	// update login stats
	now := time.Now()
	isSameDay := core.IsSameDay(now, self.user.LoginTs)

	// update login info
	self.user.LoginTs = now

	// update some timestamps
	self.hb_ts = now
	self.online_ts = now

	// sync gmails
	self.GetMailBox().SyncGMails()

	// fire
	evtmgr.Fire(gconst.Evt_LoginOnline, self, self.user.DevId, self.user.LoginIP)

	// 同一天登录不发
	if !isSameDay {
		evtmgr.Fire(gconst.Evt_PlrDailyOnline, self)
	}

	// send user-info
	self.send_user_info()

	self.UnNew()
}

func (self *Player) OnOffline(shutdown bool) {
	if !shutdown {
		log.Info(self.user.Id, self.user.Name, "is offline")
	}

	// accumulate online dur
	self.AccOnlineDur()

	// update offline ts
	self.user.Offline_ts = time.Now()

	// fire
	if !shutdown {
		evtmgr.Fire(gconst.Evt_LoginOffline, self)
	}
}

func (self *Player) IsOnline() bool {
	return self.sid != 0
}

func (self *Player) Sid() uint64 {
	return self.sid
}

func (self *Player) DB() *db.Database {
	return self.user.db
}

func (self *Player) User() *User {
	return self.user
}

func (self *Player) SendMsg(message msg.Message) {
	if self.IsOnline() {
		NetMgr.Send2Player(self.sid, message)
	}
}

func (self *Player) Logout() {
	sid := self.sid

	PlayerMgr.SetOffline(self, false)
	NetMgr.Send2Gate(sid_to_gateid(sid), &msg.GS_Kick{Sid: sid})
}

func (self *Player) HeartbeatUpdate() {
	self.hb_ts = time.Now()
}

func (self *Player) check_heartbeat() {
	if !self.IsOnline() {
		return
	}

	if time.Since(self.hb_ts).Seconds() > 60 {
		self.Logout()
	}
}

func (self *Player) update_atkpwr() {
	// update atk power
	self.user.atkpwr = self.GetTeamMgr().CalcPlrAtkPwr()

	// notify
	if self.IsOnline() {
		next_tick.Once(fmt.Sprintf("plr.ap.%s", self.GetId()), func() {
			self.SendMsg(&msg.GS_PlayerUpdateAtkPwr{
				AtkPwr: self.user.atkpwr,
			})
		})
	}

	// fire
	evtmgr.Fire(gconst.Evt_PlrAtkPwr, self, self.user.atkpwr)
}

// ============================================================================
// 基础接口

func (self *Player) GetAuthId() string {
	return self.user.AuthId
}

func (self *Player) GetId() string {
	return self.user.Id
}

func (self *Player) GetName() string {
	return self.user.Name
}

func (self *Player) GetSvr0() string {
	return self.user.Svr0
}

func (self *Player) GetSdk() string {
	return self.user.Sdk
}

func (self *Player) GetModel() string {
	return self.user.Model
}

func (self *Player) GetDevId() string {
	return self.user.DevId
}

func (self *Player) GetOs() string {
	return self.user.Os
}

func (self *Player) GetOsVer() string {
	return self.user.OsVer
}

func (self *Player) GetLoginIP() string {
	return self.user.LoginIP
}

func (self *Player) GetCreateTs() time.Time {
	return self.user.CreateTs
}

func (self *Player) IsNew() bool {
	return self.user.IsNew
}

func (self *Player) UnNew() {
	self.user.IsNew = false
}

func (self *Player) GetHead() string {
	return self.user.Head
}

func (self *Player) GetHFrame() int32 {
	return self.user.HFrame
}

func (self *Player) SetHFrame(id int32) {
	self.user.HFrame = id
}

func (self *Player) GetLevel() int32 {
	return self.user.Lv
}

func (self *Player) GetExp() int32 {
	return self.user.Exp
}

func (self *Player) GetAtkPwr() int32 {
	return self.user.atkpwr
}

func (self *Player) ChangeName(name string, f func(bool)) {
	oldname := self.user.Name

	async.Push(func() {
		// update name-db
		if !dbmgr.Center_ChangeName(oldname, name) {
			loop.Push(func() {
				f(false)
			})
			return
		}

		// update center-db
		dbmgr.Share_UpdateUserName(self.GetId(), name)

		// update game-db
		err := self.DB().Update(
			dbmgr.C_tabname_user,
			self.GetId(),
			db.M{"$set": db.M{"base.name": name}},
		)
		if err != nil {
			log.Warning("Player.ChangeName() failed:", err)
		}

		// update memory
		loop.Push(func() {
			PlayerMgr.UpdatePlayerName(self, name)

			//fire
			evtmgr.Fire(gconst.Evt_PlrChangeName, self, oldname, name)

			f(true)
		})
	})
}

func (self *Player) AddExp(v int32) {
	if v <= 0 {
		return
	}

	new_lv := self.user.Lv
	self.user.Exp += v

	// consume exp
	for {
		conf := gamedata.ConfPlayerUp.Query(new_lv)
		if conf == nil || conf.Exp == 0 {
			// 满级
			self.user.Exp = 0
			break
		}

		if self.user.Exp < conf.Exp {
			break
		}

		self.user.Exp -= conf.Exp
		new_lv++
	}

	// levelup
	self.SetLevel(new_lv)
}

func (self *Player) SetLevel(lv int32) {
	if lv == self.user.Lv {
		self.SendMsg(&msg.GS_PlayerUpdateLv{
			Level: -1, // 无变化
			Exp:   self.user.Exp,
		})
		return
	}

	// set new lv
	old := self.user.Lv
	self.user.Lv = lv

	// event: levelup
	self.OnLevelup(old, lv)
}

func (self *Player) OnLevelup(old_lv, new_lv int32) {
	// update userinfo
	async.Push(
		func() {
			dbmgr.Share_UpdateUserLv(self.GetId(), new_lv)
		},
	)

	// rewards
	op := self.GetBag().NewOp(gconst.ObjFrom_PlrLvUp)
	for i := old_lv; i < new_lv; i++ {
		conf := gamedata.ConfPlayerUp.Query(i)
		if conf != nil {
			for _, v := range conf.UpReward {
				op.Inc(v.Id, v.N)
			}
		}
	}
	rwds := op.Apply().ToMsg()

	// notify
	self.SendMsg(&msg.GS_PlayerUpdateLv{
		Level:   new_lv,
		Exp:     self.user.Exp,
		Rewards: rwds,
	})

	evtmgr.Fire(gconst.Evt_PlrLv, self, new_lv, new_lv-old_lv)
}

func (self *Player) AccOnlineDur() int32 {
	if !self.IsOnline() {
		return self.user.OnlineDur
	}

	now := time.Now()

	self.user.OnlineDur += int32(now.Sub(self.online_ts).Seconds())

	self.online_ts = now

	return self.user.OnlineDur
}

func (self *Player) GetLoginTs() time.Time {
	return self.user.LoginTs
}

func (self *Player) GetOnlineDur() int32 {
	return self.user.OnlineDur
}

func (self *Player) GetOfflineTs() time.Time {
	return self.user.Offline_ts
}

func (self *Player) GetLoginSumDays() int32 {
	return self.user.LoginSumDays
}

func (self *Player) IsActive() bool {
	ts := time.Now().Add(-gconst.PLAYER_ActiveDays * 24 * time.Hour)
	return self.user.LoginTs.After(ts)
}

func (self *Player) GetAuthRet() map[string]string {
	return self.user.AuthRet
}

// ============================================================================

// 辅助补丁函数
func (self *Player) patch() {
	evtmgr.Fire(gconst.Evt_PlrLv, self, self.user.Lv, int32(0))
	evtmgr.Fire(gconst.Evt_WLevelLv, self, self.GetWLevelLvNum(), int32(0))
}

// ============================================================================
// 扩展接口

//  atkpwr不填为玩家默认战力
func (self *Player) ToMsg_SimpleInfo(atkpwr ...int32) *msg.PlayerSimpleInfo {
	ap := self.GetAtkPwr()
	if len(atkpwr) > 0 {
		ap = atkpwr[0]
	}

	ret := &msg.PlayerSimpleInfo{
		Id:     self.GetId(),
		Name:   self.GetName(),
		Lv:     self.GetLevel(),
		Exp:    self.GetExp(),
		Head:   self.GetHead(),
		HFrame: self.GetHFrame(),
		Vip:    self.GetVipLevel(),
		SvrId:  config.CurGame.Id,
		AtkPwr: ap,
		GName:  self.GetGuildName(),
	}

	h := self.GetBag().MaxPowerHero()
	if h != nil {
		ret.ShowHero = h.Id
	}

	return ret
}

// ============================================================================
// 存盘

func (self *Player) save_timer_start() {
	loop.SetTimeout(self.save_next_time(), func() {
		// check heart-beat: just to avoid very long time suspending
		self.check_heartbeat()

		// accumulate online dur
		self.AccOnlineDur()

		// offline check
		if !self.IsOnline() {
			self.GetCounter().CheckRecover(gconst.Cnt_PlayerStrength)
			self.GetWLevel().GJLoot()
		}

		// save async
		self.save_async()
		self.save_timer_start()
	})
}

func (self *Player) save_next_time() time.Time {
	// 此方案是为了减轻玩家存盘压力
	n := PlayerMgr.NumOnline()
	t := 1800
	if n > 10000 {
		t = 2400
	} else if n > 20000 {
		t = 3600
	}

	if self.IsOnline() {
		// 在线存盘间隔
		return time.Now().Add(time.Duration(t+rand.Intn(t)) * time.Second)
	} else {
		// 离线存盘间隔
		return time.Now().Add(time.Duration(t*2+rand.Intn(t*2)) * time.Second)
	}
}

// 异步存盘：用于定时存盘
func (self *Player) save_async() {
	// clone
	obj := core.CloneBsonObject(self.user)

	// async save
	async.Push(
		func() {
			err := self.DB().Update(
				dbmgr.C_tabname_user,
				self.user.Id,
				db.M{"$set": db.M{"base": obj}},
			)
			if err != nil {
				log.Error("saving player failed:", err)
			}
		},
	)
}

// 同步存盘：用于停服时存盘
func (self *Player) save() {
	err := self.DB().Update(
		dbmgr.C_tabname_user,
		self.user.Id,
		db.M{"$set": db.M{"base": self.user}},
	)
	if err != nil {
		log.Error("final saving player failed:", err)
	}
}

// ============================================================================
// 重置处理

func (self *Player) Reset_GetTime() time.Time {
	return self.user.Rst_ts
}

func (self *Player) Reset_SetTime(ts time.Time) {
	self.user.Rst_ts = ts
}

func (self *Player) Reset_Daily() {
	evtmgr.Fire(gconst.Evt_PlrResetDaily, self)
}

func (self *Player) Reset_Weekly() {
	evtmgr.Fire(gconst.Evt_PlrResetWeekly, self)
}

func (self *Player) Reset_Monthly() {
	evtmgr.Fire(gconst.Evt_PlrResetMonthly, self)
}

// ============================================================================
// mods

// ============================================================================
// 封号相关

func (self *Player) GetBanTs() time.Time {
	return self.user.BanTs
}

func (self *Player) IsBan() bool {
	return time.Now().Before(self.user.BanTs)
}

// ============================================================================
// 背包

func (self *Player) GetBag() *comp.Bag {
	return self.user.Bag
}

// ============================================================================
// 计数器相关

func (self *Player) GetCounter() *comp.Counter {
	return self.user.Counter
}

// ============================================================================
// 充值

func (self *Player) GetBill() *bill.Bill {
	return self.user.Bill
}

// 总充值基准货币数量
func (self *Player) GetBillTotalBaseCcy() int64 {
	return self.user.Bill.TotalBaseCcy
}

// ============================================================================
// 条件统计

func (self *Player) GetAttainTab() *attaintab.AttainTab {
	return self.user.AttainTab
}

func (self *Player) GetAttainObjVal(oid int32) float64 {
	return self.GetAttainTab().GetObjVal(oid)
}

// ============================================================================
// 邮箱

func (self *Player) GetMailBox() *mail.MailBox {
	return self.user.MailBox
}

func (self *Player) CheckGmailDeliverCond(cond string) bool {
	// empty cond: allow delivery
	if strings.TrimSpace(cond) == "" {
		return true
	}

	// check conds
	re := regexp.MustCompile(`^\s*([\w_]+)\s*([<=>!]{1,2})(.+)$`)

	for _, e := range strings.Split(cond, "|") {
		arr := re.FindStringSubmatch(e)
		if arr == nil {
			return false
		}

		name := arr[1]
		op := arr[2]
		val := strings.TrimSpace(arr[3])

		var b bool

		switch name {
		case "lv":
			b = compare_int32(self.GetLevel(), core.Atoi32(val), op)

		case "cdate":
			b = compare_date(self.user.CreateTs, core.ParseTime(val), op)

		case "ldate":
			b = compare_date(self.user.LoginTs, core.ParseTime(val), op)

		default:
			b = false
		}

		if !b {
			return false
		}
	}

	return true
}

func compare_int32(left, right int32, op string) bool {
	switch op {
	case "<":
		return left < right
	case "<=":
		return left <= right
	case ">":
		return left > right
	case ">=":
		return left >= right
	case "==":
		return left == right
	case "!=":
		return left != right
	}

	return false
}

func compare_date(left, right time.Time, op string) bool {
	switch op {
	case "<":
		return left.Before(right)
	case "<=":
		return !left.After(right)
	case ">":
		return left.After(right)
	case ">=":
		return !left.Before(right)
	case "==":
		return left.Equal(right)
	case "!=":
		return !left.Equal(right)
	}

	return false
}

// ============================================================================
// 公会

func (self *Player) GetGuild() *guild.Guild {
	return self.user.Guild
}

func (self *Player) GetGuildId() string {
	return self.user.GuildId
}

func (self *Player) GetGuildName() string {
	if self.user.Guild == nil {
		return ""
	} else {
		return self.user.Guild.Name
	}
}

func (self *Player) GetGuildIcon() int32 {
	if self.user.Guild == nil {
		return 1101
	} else {
		return self.user.Guild.Icon
	}
}

func (self *Player) GetGuildRank() int32 {
	if self.user.Guild == nil {
		return 1101
	} else {
		m := self.user.Guild.FindMember(self.GetId())
		if m == nil {
			return 0
		} else {
			return m.Rank
		}
	}
}

func (self *Player) BindGuild(gldid string) {
	// bind
	if gldid != "" {
		gld := guild.GuildMgr.FindGuild(gldid)
		if gld != nil {
			m := gld.FindMember(self.GetId())
			if m != nil {
				self.user.GuildId = gld.Id
				self.user.Guild = gld
				return
			}
		}
	}

	// failed
	self.user.GuildId = ""
	self.user.Guild = nil
}

func (self *Player) GetGuildPlrData() *guild.GuildPlrData {
	return self.user.GuildPlrData
}

func (self *Player) IsSameGuild(plrid string) bool {
	g := self.GetGuild()
	return g != nil && g.Members[plrid] != nil
}

// ============================================================================
// 模块开启

func (self *Player) GetMOpen() *mopen.MOpen {
	return self.user.MOpen
}

func (self *Player) IsModuleOpen(mid int32) bool {
	return self.user.MOpen.IsOpen(mid)
}

// ============================================================================
// 新手

func (self *Player) GetTutorial() *tutorial.Tutorial {
	return self.user.Tutorial
}

// ============================================================================
// 云数据

func (self *Player) GetCloud() *cloud.Cloud {
	return self.user.Cloud
}

// ============================================================================
// 杂项数据

func (self *Player) GetMisc() *misc.Misc {
	return self.user.Misc
}

// ============================================================================
// 活动

func (self *Player) GetAct() *act.Act {
	return self.user.Act
}

func (self *Player) GetActRawData(name string) interface{} {
	return self.user.Act.GetActRawData(name)
}

// ============================================================================
// 日常任务

func (self *Player) GetTaskDaily() *taskdaily.TaskDaily {
	return self.user.TaskDaily
}

// ============================================================================
// 成就任务

func (self *Player) GetTaskAchv() *taskachv.TaskAchv {
	return self.user.TaskAchv
}

// ============================================================================
// 抽卡

func (self *Player) GetDraw() *draw.Draw {
	return self.user.Draw
}

// ============================================================================
// vip
func (self *Player) GetVip() *vip.Vip {
	return self.user.Vip
}

func (self *Player) GetVipLevel() int32 {
	return self.GetVip().Lv
}

func (self *Player) AddVipExp(v int32) {
	self.GetVip().AddExp(v)
}

// ============================================================================
// battle

func (self *Player) IsTeamFormationValid(tf *msg.TeamFormation) bool {
	if tf == nil {
		return false
	}

	bag := self.GetBag()

	// check heroes
	m := make(map[int32]bool) // [pos]

	for seq, pos := range tf.Formation {
		// pos MUST be valid and unique
		if pos < 0 || pos >= gconst.Team_MaxSlots || m[pos] {
			return false
		}

		// we MUST own the hero
		hero := bag.FindHero(seq)
		if hero == nil {
			return false
		}

		m[pos] = true
	}

	// at least 1 hero is needed
	if len(m) == 0 {
		return false
	}

	// ok
	return true
}

func (self *Player) GetTeamAtkPwr(tf *msg.TeamFormation) int32 {
	ap := int32(0)
	if tf != nil {
		for seq := range tf.Formation {
			hero := self.GetBag().FindHero(seq)
			if hero != nil {
				ap += hero.GetAtkPower()
			}
		}
	}

	return ap
}

// teamAtkPwr true为本队战力, false或不填为玩家默认战力
func (self *Player) ToMsg_BattleTeam(tf *msg.TeamFormation, teamAtkPwr ...bool) *msg.BattleTeam {
	if tf == nil {
		return nil
	}

	bag := self.GetBag()

	// fighters
	ap := int32(0)
	fts := make([]*msg.BattleFighter, 0, len(tf.Formation))
	for seq, pos := range tf.Formation {
		hero := bag.FindHero(seq)
		if hero == nil {
			continue
		}

		fts = append(fts, hero.ToMsg_BattleFighter(pos))

		ap += hero.GetAtkPower()
	}

	if !core.DefFalse(teamAtkPwr) {
		ap = self.GetAtkPwr()
	}

	// ok
	return &msg.BattleTeam{
		Player:   self.ToMsg_SimpleInfo(ap),
		Fighters: fts,
	}
}

// ============================================================================
// 推图

func (self *Player) GetWLevel() *wlevel.WLevel {
	return self.user.WLevel
}

func (self *Player) GetWLevelLvNum() int32 {
	return self.GetWLevel().LvNum
}

// ============================================================================
// 酒馆委派

func (self *Player) GetAppoint() *appoint.Appoint {
	return self.user.Appoint
}

// ============================================================================
// 爬塔

func (self *Player) GetTower() *tower.Tower {
	return self.user.Tower
}

func (self *Player) GetTowerLvNum() int32 {
	return self.GetTower().LvNum
}

// ============================================================================
// 阵容管理

func (self *Player) GetTeamMgr() *teammgr.TeamMgr {
	return self.user.TeamMgr
}

func (self *Player) IsSetTeam(tp int32) bool {
	return self.GetTeamMgr().IsSetTeam(tp)
}

func (self *Player) GetTeam(tp int32) *msg.TeamFormation {
	return self.GetTeamMgr().GetTeam(tp)
}

func (self *Player) SetTeam(tp int32, T *msg.TeamFormation) {
	self.GetTeamMgr().SetTeam(tp, T)
}

// ============================================================================
// 竞技场

func (self *Player) GetArena() *arena.Arena {
	return self.user.Arena
}

func (self *Player) GetArenaScore() int32 {
	return self.GetArena().GetScore()
}

func (self *Player) AddArenaScore(v int32) {
	self.GetArena().AddScore(v)
}

func (self *Player) AddArenaReplay(rid string, sscore, escore, ascore int32, plrInfo *msg.PlayerSimpleInfo) {
	self.GetArena().AddArenaReplay(rid, sscore, escore, ascore, plrInfo)
}

// ============================================================================
// 头像框

func (self *Player) GetHFrameStore() *hframestore.HFrameStore {
	return self.user.HFrameStore
}

// ============================================================================
// 商店

func (self *Player) GetShop() *shop.Shop {
	return self.user.Shop
}

// ============================================================================
// 奇迹之盘

func (self *Player) GetMarvelRoll() *marvelroll.MarvelRoll {
	return self.user.MarvelRoll
}

// ============================================================================
// 好友

func (self *Player) GetFriend() *friend.Friend {
	return self.user.Friend
}

func (self *Player) IsFriend(plrid string) bool {
	return self.GetFriend().IsFriend(plrid)
}

// ============================================================================
// 远征

func (self *Player) GetCrusade() *crusade.Crusade {
	return self.user.Crusade
}

// ============================================================================
// 月票

func (self *Player) GetMonthTicket() *monthticket.MonthTicket {
	return self.user.MonthTicket
}

// ============================================================================
// 推送礼包

func (self *Player) GetPushGift() *pushgift.PushGift {
	return self.user.PushGift
}

// ============================================================================
// 礼包商店

func (self *Player) GetGiftShop() *giftshop.GiftShop {
	return self.user.GiftShop
}

// ============================================================================
// 特权卡

func (self *Player) GetPrivCard() *privcard.PrivCard {
	return self.user.PrivCard
}

func (self *Player) IsPrivCardValid(id int32) bool {
	return self.GetPrivCard().IsPrivCardValid(id)
}

func (self *Player) AddPrivCard(id int32) {
	self.GetPrivCard().AddPrivCard(id)
}

// ============================================================================
// 签到
func (self *Player) GetSignDaily() *signdaily.SignDaily {
	return self.user.SignDaily
}

// ============================================================================
// 每月任务

func (self *Player) GetTaskMonth() *taskmonth.TaskMonth {
	return self.user.TaskMonth
}

// ============================================================================
// 七日之约

func (self *Player) GetDaySign() *daysign.DaySign {
	return self.user.DaySign
}

// ============================================================================
// 榜单玩法

func (self *Player) GetRankPlay() *rank.RankPlay {
	return self.user.RankPlay
}

// ============================================================================
// 开服庆典(七日目标)

func (self *Player) GetTargetDays() *targetdays.TargetDays {
	return self.user.TargetDays
}

// ============================================================================
// 进阶之路

func (self *Player) GetTaskGrow() *taskgrow.TaskGrow {
	return self.user.TaskGrow
}

// ============================================================================
// 推图基金

func (self *Player) GetWLevelFund() *wlevelfund.WLevelFund {
	return self.user.WLevelFund
}

// ============================================================================
// 成长基金

func (self *Player) GetGrowFund() *growfund.GrowFund {
	return self.user.GrowFund
}

// ============================================================================
// 超值首充

func (self *Player) GetBillFirst() *billfirst.BillFirst {
	return self.user.BillFirst
}

// ============================================================================
// 公会战

func (self *Player) GetGWar() *gwar.GWar {
	return self.user.GWar
}

// ============================================================================
// 裂隙怪物

func (self *Player) GetRift() *rift.Rift {
	return self.user.Rift
}

// ============================================================================
// 英雄皮肤

func (self *Player) GetHeroSkin() *heroskin.HeroSkin {
	return self.user.HeroSkin
}

// ============================================================================
// 天梯

func (self *Player) GetLadder() *ladder.Ladder {
	return self.user.Ladder
}

// ============================================================================
// 推图十连

func (self *Player) GetWLevelDraw() *wleveldraw.WLevelDraw {
	return self.user.WLevelDraw
}

// ============================================================================
// 本服杯赛

func (self *Player) GetWarCup() *warcup.WarCup {
	return self.user.WarCup
}

// ============================================================================
// wboss

func (self *Player) GetWBoss() *wboss.WBoss {
	return self.user.WBoss
}

// ============================================================================
// invite

func (self *Player) GetInvite() *invite.Invite {
	return self.user.Invite
}
