package main

import (
	"flag"
	"fw/src/core"
	"fw/src/core/log"
	"fw/src/router/app"
	"fw/src/router/handler"
	"fw/src/shared/config"
	"math/rand"
	"os"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	// parse command line
	arg_file := flag.String("config", "config.json", "config file")
	arg_server := flag.String("server", "router1", "server name")
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
		} else {
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
	// start net mgr
	handler.Init()
	app.NetMgr.Start()

	// app started
	log.Notice("router started")
}

func stop() {
	// stop net mgr
	app.NetMgr.Stop()

	// app stopped
	log.Notice("router stopped")
}
