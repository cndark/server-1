package app

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/game/app/comp"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/appoint"
	"fw/src/game/app/modules/arena"
	"fw/src/game/app/modules/attaintab"
	"fw/src/game/app/modules/bill"
	"fw/src/game/app/modules/billfirst"
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
	"time"
)

// ============================================================================

type User struct {
	Id       string    `bson:"-"`         // Id
	AuthId   string    `bson:"authid"`    // 认证 Id
	Svr0     string    `bson:"svr0"`      // 原始服名称
	Svr      string    `bson:"svr"`       // 当前服名称
	Sdk      string    `bson:"sdk"`       // sdk
	Model    string    `bson:"model"`     // 登录设备型号
	DevId    string    `bson:"devid"`     // 登录设备码
	Os       string    `bson:"os"`        // 登录操作系统
	OsVer    string    `bson:"osver"`     // 登录操作系统版本
	CreateTs time.Time `bson:"create_ts"` // 创建时间
	LoginTs  time.Time `bson:"login_ts"`  // 上次登录时间
	LoginIP  string    `bson:"login_ip"`  // 上次登录 IP

	LoginCtDays  int32     `bson:"login_ctdays"`  // 最近连续登录天数
	LoginSumDays int32     `bson:"login_sumdays"` // 累计登陆天数
	OnlineDur    int32     `bson:"online_dur"`    // 累计在线时长 (秒)
	Offline_ts   time.Time `bson:"off_ts"`        // 上次离线时间
	Rst_ts       time.Time `bson:"rst_ts"`        // 上次重置时间
	BanTs        time.Time `bson:"ban_ts"`        // 封号时间戳
	IsNew        bool      `bson:"isnew"`         // 是否是新号

	Name         string                   `bson:"name"`        // 名字
	Head         string                   `bson:"head"`        // 头像
	HFrame       int32                    `bson:"hframe"`      // 头像框
	Lv           int32                    `bson:"lv"`          // 等级
	Exp          int32                    `bson:"exp"`         // 经验
	Bag          *comp.Bag                `bson:"bag"`         // 背包
	Counter      *comp.Counter            `bson:"counter"`     // 计数
	Bill         *bill.Bill               `bson:"bill"`        // 充值
	AttainTab    *attaintab.AttainTab     `bson:"attaintab"`   // 条件统计表数据
	MailBox      *mail.MailBox            `bson:"mailbox"`     // 邮箱
	GuildPlrData *guild.GuildPlrData      `bson:"gld_plrdata"` // 家族跟着玩家数据
	MOpen        *mopen.MOpen             `bson:"mopen"`       // 模块开启
	Tutorial     *tutorial.Tutorial       `bson:"tutorial"`    // 新手
	Cloud        *cloud.Cloud             `bson:"cloud"`       // 云数据
	Misc         *misc.Misc               `bson:"misc"`        // 杂项数据
	Act          *act.Act                 `bson:"act"`         // 活动
	TaskDaily    *taskdaily.TaskDaily     `bson:"taskdaily"`   // 日常任务
	TaskAchv     *taskachv.TaskAchv       `bson:"taskachv"`    // 成就任务
	Draw         *draw.Draw               `bson:"draw"`        // 抽卡
	WLevel       *wlevel.WLevel           `bson:"wlevel"`      // 推图
	Appoint      *appoint.Appoint         `bson:"appoint"`     // 酒馆委派
	Tower        *tower.Tower             `bson:"tower"`       // 爬塔
	TeamMgr      *teammgr.TeamMgr         `bson:"teammgr"`     // 阵容管理
	Arena        *arena.Arena             `bson:"arena"`       // 竞技场
	HFrameStore  *hframestore.HFrameStore `bson:"hframestore"` // 竞技场
	Shop         *shop.Shop               `bson:"shop"`        // 商店
	MarvelRoll   *marvelroll.MarvelRoll   `bson:"marvelroll"`  // 奇迹之盘
	Friend       *friend.Friend           `bson:"friend"`      // 好友
	Crusade      *crusade.Crusade         `bson:"crusade"`     // 远征
	Vip          *vip.Vip                 `bson:"vip"`         // vip
	MonthTicket  *monthticket.MonthTicket `bson:"monthticket"` // 月票
	PushGift     *pushgift.PushGift       `bson:"pushgift"`    // 推送礼包
	GiftShop     *giftshop.GiftShop       `bson:"giftshop"`    // 礼包商店
	PrivCard     *privcard.PrivCard       `bson:"privcard"`    // 特权卡
	SignDaily    *signdaily.SignDaily     `bson:"signdaily"`   // 每日签到
	TaskMonth    *taskmonth.TaskMonth     `bson:"taskmonth"`   // 每月任务
	DaySign      *daysign.DaySign         `bson:"daysign"`     // 七日之约
	RankPlay     *rank.RankPlay           `bson:"rankplay"`    // 榜单玩法
	TargetDays   *targetdays.TargetDays   `bson:"targetdays"`  // 开服庆典(七日目标)
	TaskGrow     *taskgrow.TaskGrow       `bson:"taskgrow"`    // 进阶之路
	WLevelFund   *wlevelfund.WLevelFund   `bson:"wlevelfund"`  // 推图基金
	GrowFund     *growfund.GrowFund       `bson:"growfund"`    // 推图基金
	BillFirst    *billfirst.BillFirst     `bson:"billfirst"`   // 超值首充
	GWar         *gwar.GWar               `bson:"gwar"`        // 公会战
	Rift         *rift.Rift               `bson:"riftmonster"` // 裂隙数据(注意bson;db名字)
	HeroSkin     *heroskin.HeroSkin       `bson:"heroskin"`    // 英雄皮肤库
	Ladder       *ladder.Ladder           `bson:"ladder"`      // 天梯
	WLevelDraw   *wleveldraw.WLevelDraw   `bson:"wleveldraw"`  // 推图十连
	WarCup       *warcup.WarCup           `bson:"warcup"`      // 本服杯赛
	WBoss        *wboss.WBoss             `bson:"wboss"`       // 世界boss
	Invite       *invite.Invite           `bson:"invite"`      // 收藏分享

	GuildId string       `bson:"guildid"` // 公会 Id
	Guild   *guild.Guild `bson:"-"`       // 公会

	atkpwr int32 `bson:"-"` // 战力

	AuthRet map[string]string `bson:"-"` // 认证信息(不存盘)
	db      *db.Database      `bson:"-"` // db
}

// ============================================================================

func create_user(uid string, f func(*User)) *User {
	// get user db
	dbname := get_user_dbname(uid)
	udb := dbmgr.UserDB(dbname)
	if udb == nil {
		log.Error("get user db failed:", dbname)
		log.Error(core.Callstack())
		return nil
	}

	// new user
	user := &User{}

	// callback
	if f != nil {
		f(user)
	}

	// init data
	user.Id = uid
	user.CreateTs = time.Now()

	user.LoginCtDays = 0
	user.LoginSumDays = 0
	user.OnlineDur = 0
	user.Offline_ts = time.Unix(0, 0)
	user.Rst_ts = time.Unix(0, 0)
	user.BanTs = time.Unix(0, 0)
	user.IsNew = true

	user.Head = user.AuthId
	user.HFrame = 0
	user.Lv = 1
	user.Bag = comp.NewBag()
	user.Counter = comp.NewCounter()
	user.Bill = bill.NewBill()
	user.AttainTab = attaintab.NewAttainTab()
	user.MailBox = mail.NewMailBox()
	user.GuildPlrData = guild.NewGuildPlrData()
	user.MOpen = mopen.NewMOpen()
	user.Tutorial = tutorial.NewTutorial()
	user.Cloud = cloud.NewCloud()
	user.Misc = misc.NewMisc()
	user.Act = act.NewAct()
	user.TaskDaily = taskdaily.NewTaskDaily()
	user.TaskAchv = taskachv.NewTaskAchv()
	user.Draw = draw.NewDraw()
	user.WLevel = wlevel.NewWLevel()
	user.Appoint = appoint.NewAppoint()
	user.Tower = tower.NewTower()
	user.TeamMgr = teammgr.NewTeamMgr()
	user.Arena = arena.NewArena()
	user.HFrameStore = hframestore.NewHFrameStore()
	user.Shop = shop.NewShop()
	user.MarvelRoll = marvelroll.NewMarvelRoll()
	user.Friend = friend.NewFriend()
	user.Crusade = crusade.NewCrusade()
	user.Vip = vip.NewVip()
	user.MonthTicket = monthticket.NewMonthTicket()
	user.PushGift = pushgift.NewPushGift()
	user.GiftShop = giftshop.NewGiftShop()
	user.PrivCard = privcard.NewPrivCard()
	user.SignDaily = signdaily.NewSignDaily()
	user.TaskMonth = taskmonth.NewTaskMonth()
	user.DaySign = daysign.NewDaySign()
	user.RankPlay = rank.NewRankPlay()
	user.TargetDays = targetdays.NewTargetDays()
	user.TaskGrow = taskgrow.NewTaskGrow()
	user.WLevelFund = wlevelfund.NewWLevelFund()
	user.GrowFund = growfund.NewGrowFund()
	user.BillFirst = billfirst.NewBillFirst()
	user.GWar = gwar.NewGWar()
	user.Rift = rift.NewRift()
	user.HeroSkin = heroskin.NewHeroSkin()
	user.Ladder = ladder.NewLadder()
	user.WLevelDraw = wleveldraw.NewWLevelDraw()
	user.WarCup = warcup.NewWarCup()
	user.WBoss = wboss.NewWBoss()
	user.Invite = invite.NewInvite()

	user.GuildId = ""
	user.Guild = nil

	// --------------------------------

	// save to db
	err := udb.Insert(dbmgr.C_tabname_user, db.M{"_id": user.Id, "base": user})
	if err != nil {
		log.Error("create user failed:", uid, err)
		return nil
	}

	// update user name into center-db
	dbmgr.Share_UpdateUserName(user.Id, user.Name)

	// bind
	user.db = udb

	// return
	return user
}
