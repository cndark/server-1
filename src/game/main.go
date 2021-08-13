package main

import (
	"flag"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/core/sched/resetter"
	"fw/src/game/app"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/glog"
	"fw/src/game/app/gsvc"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/calendar"
	"fw/src/game/app/modules/chat"
	"fw/src/game/app/modules/guild"
	"fw/src/game/app/modules/lamp"
	"fw/src/game/app/modules/loader"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/mdata"
	"fw/src/game/app/modules/rank"
	"fw/src/game/app/modules/refund"
	"fw/src/game/app/modules/statustab"
	"fw/src/game/app/modules/svrgrp"
	"fw/src/game/app/modules/utils"
	"fw/src/game/app/modules/worlddata"
	"fw/src/game/app/stats"
	"fw/src/game/handler"
	"fw/src/game/perfmon"
	"fw/src/shared/config"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	// parse command line
	arg_file := flag.String("config", "config.json", "config file")
	arg_server := flag.String("server", "game1", "server name")
	arg_log := flag.String("log", "", "log file")

	flag.Parse()

	// open log
	log.Open(*arg_log)

	// load config
	config.Parse(*arg_file, *arg_server)

	// pprof
	if config.CurGame.Id == 1 {
		go func() {
			runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
			runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
			http.ListenAndServe(":8011", nil)
		}()
	}

	// signal
	quit := make(chan int)
	core.OnSignal(func(s os.Signal) {
		if s == syscall.SIGHUP {
			log.Warning("reloading ...")

			loop.Push(func() {
				config.Parse(*arg_file, *arg_server)

				dbmgr.Open()
				evtmgr.Fire(gconst.Evt_ConfReload)
			})
		} else if s == syscall.SIGINT || s == syscall.SIGTERM {
			log.Warning("shutdown signal received ...")
			close(quit)
		}
	})

	start()
	<-quit
	stop()

	// close log
	log.Close()
}

func start() {
	// ============================================================================
	// INIT WORLD
	// ============================================================================

	// open db mgr
	dbmgr.Open()

	// load gamedata
	gamedata.Load()

	// init async-op
	async.Init([]int{
		3, // default
		1, // mail
		1, // bill
		1, // stats
		3, // glog
	})

	// init utils
	utils.Init(app.PlayerMgr, app.NetMgr, statustab.StatusTab)

	// loader
	loader.LoadModules()

	// worlddata
	worlddata.Init()

	// svrgrp
	svrgrp.Init()

	// act
	act.Open()

	// rank
	rank.Init()

	// gmail box
	mail.GMailBox.Init()

	// refund
	refund.Open()

	// guild mgr
	guild.GuildMgr.Open()

	// module data
	mdata.Open()

	// battle
	battle.Init()

	chat.Open()

	lamp.Init()

	// ============================================================================
	// WORLD READY
	// ============================================================================

	evtmgr.Fire(gconst.Evt_WorldReady)

	// ============================================================================
	// LOADING PLAYERS
	// ============================================================================

	app.PlayerMgr.Open()

	// ============================================================================
	// START WORLD
	// ============================================================================

	// resetter
	resetter.Start()

	// calendar
	calendar.Start()

	// ============================================================================
	// AUX
	// ============================================================================

	// start gservice
	gsvc.Start()

	// stats & glog
	stats.Start()
	glog.Start()

	// ============================================================================
	// RUN LOOP
	// ============================================================================

	// start net mgr
	handler.Init()
	app.NetMgr.Start()

	// start async & loop
	async.Start()
	loop.Run()

	// perfmon
	perfmon.Start()

	// fire
	evtmgr.Fire(gconst.Evt_ServerStart)

	// app started
	log.Notice("game started")
}

func stop() {
	// ============================================================================
	// STOP WORLD
	// ============================================================================

	// stop net mgr
	app.NetMgr.Stop()

	// stop loop & async
	log.Info("stopping loop ...")
	loop.Stop()

	log.Info("stopping async ...")
	async.Stop()

	// offline all players
	app.PlayerMgr.OfflinePlayers(0)

	// ============================================================================
	// FINAL SAVE MODULES
	// ============================================================================

	chat.Close()
	mdata.Close()
	act.Close()
	worlddata.Close()

	// ============================================================================
	// FINAL SAVE GUILD & PLAYER
	// ============================================================================

	app.PlayerMgr.Close()
	guild.GuildMgr.Close()

	// close db mgr
	dbmgr.Close()

	// fire
	evtmgr.Fire(gconst.Evt_ServerStop)

	// app stopped
	log.Notice("game stopped")
}
