package gconst

// ============================================================================
// 事件

const (
	// ==============================================
	// global

	Evt_ServerStart   = "svr.start"
	Evt_ServerStop    = "svr.stop"
	Evt_ServerNewUser = "svr.newuser"

	Evt_WorldReady     = "world.ready"
	Evt_SvrGrpReady    = "svrgrp.ready"
	Evt_OnlineNum      = "online.num"
	Evt_ConfReload     = "conf.reload"
	Evt_GameDataLoaded = "gamedata.loaded"
	Evt_GameDataReload = "gamedata.reload"

	Evt_GlobalResetDaily   = "global.rst.daily"
	Evt_GlobalResetWeekly  = "global.rst.weekly"
	Evt_GlobalResetMonthly = "global.rst.monthly"

	// ==============================================
	// gs push pull

	Evt_GsPull_PlrInfo         = "gpl.plrinfo"
	Evt_GsPush_RankCacheExpire = "gps.rank.cache.expire"

	Evt_GsPush_ChatCross = "gps.chat.cross"
	Evt_GsPush_LampCross = "gps.lamp.cross"

	Evt_GsPush_GWarGldJf = "gps.gwar.gldjf"

	Evt_GsPull_Ladder_AddPlrRequest = "gpl.ldr.addplr.req"
	Evt_GsPush_Ladder_AddPlr        = "gps.ldr.addplr"
	Evt_GsPull_Ladder_FightLock     = "gpl.ldr.ft.lock"
	Evt_GsPush_Ladder_FightUnlock   = "gps.ldr.ft.unlock"
	Evt_GsPull_Ladder_Fight         = "gpl.ldr.ft"
	Evt_GsPush_Ladder_SyncPos       = "gps.ldr.sync.pos"

	// ==============================================
	// module

	Evt_PlrLv          = "plr.lv"
	Evt_PlrChangeName  = "plr.chgname"
	Evt_PlrAtkPwr      = "plr.atkpwr"
	Evt_PlrDailyOnline = "plr.daily.online"

	Evt_PlrResetDaily   = "plr.rst.daily"
	Evt_PlrResetWeekly  = "plr.rst.weekly"
	Evt_PlrResetMonthly = "plr.rst.monthly"

	Evt_LoginOnline  = "login.online"
	Evt_LoginOffline = "login.offline"

	Evt_BillGen        = "bill.gen"
	Evt_BillDone       = "bill.done"
	Evt_BillStats      = "bill.stats"
	Evt_BillGrowFund   = "bill.grow.fund"
	Evt_BillWLevelFund = "bill.wlevel.fund"
	Evt_BillFirst      = "bill.first"

	Evt_BagChg        = "bag.chg"
	Evt_FieldChange   = "fld.chg"
	Evt_CcyAdd        = "ccy.add"
	Evt_CcyDel        = "ccy.del"
	Evt_ItemAdd       = "item.add"
	Evt_ItemDel       = "item.del"
	Evt_HeroAdd       = "hero.add"
	Evt_HeroDel       = "hero.del"
	Evt_HeroLv        = "hero.lv"
	Evt_HeroStar      = "hero.star"
	Evt_HeroReset     = "hero.rst"
	Evt_HeroDecompose = "hero.decompose"
	Evt_HeroAtkPower  = "hero.atkpwr"
	Evt_HeroTrinketLv = "hero.trinketlv"
	Evt_RelicAdd      = "relic.add"

	Evt_ArmorCompose = "ar.compose"

	Evt_ModsHero = "mods.hero"

	Evt_MOpen = "module.open"

	Evt_GuildCreate     = "gld.create"
	Evt_GuildDestroy    = "gld.destroy"
	Evt_GuildJoin       = "gld.join"
	Evt_GuildLeave      = "gld.leave"
	Evt_GuildMemberRank = "gld.mrk"
	Evt_GuildLv         = "gld.lv"
	Evt_GuildExpAdd     = "gld.expadd"
	Evt_GuildExpPlrAdd  = "gld.expplradd"
	Evt_GuildChange     = "gld.chg"
	Evt_GuildDonate     = "gld.donate"
	Evt_GuildUserChange = "gld.userchg"
	Evt_GuildBossFight  = "gld.bossft"
	Evt_GuildSign       = "gld.sign"
	Evt_GuildWish       = "gld.wish"
	Evt_GuildHelp       = "gld.help"
	Evt_GuildOrderClose = "gld.orderclose"

	Evt_Tutorial = "tutorial"

	Evt_TaskDaily_Take  = "taskdaily.take"
	Evt_TaskDaily_Fin   = "taskdaily.fin"
	Evt_TaskAchv_Take   = "taskachv.take"
	Evt_TargetDays_Take = "targetdays.take"

	Evt_TaskMonth_Fin = "taskmonth.fill"

	Evt_TaskGrow_Take = "taskgrow.take"

	Evt_Draw = "draw"

	Evt_WLevelLv       = "wlevel.lv"
	Evt_WLevelFight    = "wlevel.ft"
	Evt_WLevelGj       = "wlevel.gj"
	Evt_WLevelOneKeyGj = "wlevel.onekeygj"

	Evt_TowerLv    = "tower.lv"
	Evt_TowerFight = "tower.ft"
	Evt_TowerRaid  = "tower.raid"

	Evt_DfdTeamUpdate = "dfdteam.update"

	Evt_ShopBuy = "shop.buy"

	Evt_VipLv = "vip.lv"

	Evt_SendChat = "send.chat"

	Evt_PrivCardAdd = "privcard.add"
	Evt_SignDaily   = "sign.daily"

	Evt_FriendAdd   = "friend.add"
	Evt_FriendGive  = "friend.give"
	Evt_AppointFin  = "appoint.fin"
	Evt_AppointSend = "appoint.send"

	Evt_MarvelRoll = "marvelroll"

	Evt_GoldHand = "gold.hand"

	Evt_ArenaFight = "arena.ft"
	Evt_ArenaScore = "arena.score"
	Evt_ArenaRank  = "arena.rank"

	Evt_CrusadeFight = "crusade.ft"

	Evt_RankLike = "rank.like"

	Evt_PushGift = "pushgift"

	Evt_SharedGame = "shared.game"

	Evt_RiftMonsterFight = "rift.monster.ft"
	Evt_RiftMineOccupy   = "rift.mine.occ"
	Evt_RiftMineTake     = "rift.mine.take"
	Evt_RiftBoxTake      = "rift.box.take"

	Evt_ActGift = "act.gift"

	Evt_LadderFight = "ladder.fight"

	Evt_WarCupGuess    = "warcup.guess"
	Evt_WarCupGuessWin = "warcup.guess.win"
	Evt_WarCupChat     = "warcup.chat"
	Evt_WarCupWatch    = "warcup.watch"

	Evt_OfflineCounterFull  = "counter.full"
	Evt_OfflineWlevelGjFull = "wlevel.gjfull"

	Evt_ActMaze_Click = "actmaze.click"
	Evt_ActMaze_Event = "actmaze.event"
	Evt_ActMaze_Score = "actmaze.score"
)
