package main

import (
	"flag"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/gconst"
	"fw/src/gate/app"
	"fw/src/gate/app/dbmgr"
	"fw/src/gate/handler"
	"fw/src/shared/config"
	"os"
	"syscall"
)

func main() {
	// parse command line
	arg_file := flag.String("config", "config.json", "config file")
	arg_server := flag.String("server", "gate1", "server name")
	arg_log := flag.String("log", "", "log file")

	flag.Parse()

	// open log
	log.Open(*arg_log)

	// load config
	config.Parse(*arg_file, *arg_server)

	// signal
	quit := make(chan int)
	core.OnSignal(func(s os.Signal) {
		if s == syscall.SIGHUP {
			log.Warning("reloading ...")

			config.Parse(*arg_file, *arg_server)
			app.NetMgr.ModifyGameConnections()
		} else if s == syscall.Signal(0xa) {
			log.Warning("kickgame ...")

			app.NetMgr.KickGamesPlayer()
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
	// open db mgr
	dbmgr.Open()

	// init session
	app.InitSession()

	// start net mgr
	handler.Init()
	app.NetMgr.Start()

	// fire
	evtmgr.Fire(gconst.Evt_ServerStart)

	// app started
	log.Notice("gate started")
}

func stop() {
	// stop net mgr
	app.NetMgr.Stop()

	// close db mgr
	dbmgr.Close()

	// fire
	evtmgr.Fire(gconst.Evt_ServerStop)

	// app stopped
	log.Notice("gate stopped")
}
