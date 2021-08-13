package main

import (
	"flag"
	"fmt"
	"fw/src/bot/app"
	"fw/src/bot/app/job"
	"fw/src/bot/botconf"
	"fw/src/bot/handler"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/gamedata"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	// parse flags
	conf_file := flag.String("config", "bot.json", "bot config file")
	flag.Parse()

	if *conf_file == "" {
		flag.PrintDefaults()
		return
	}

	// load config
	botconf.Parse(*conf_file)

	// open log
	log.Open("")

	// signal
	quit := make(chan int, 1)
	core.OnSignal(func(s os.Signal) {
		log.Warning("shutdown signal received ...")
		quit <- 0
	})

	evtmgr.On("all.bot.dies", func(...interface{}) {
		log.Warning("ALL bots dies ...")
		quit <- 0
	})

	// start pprof
	if botconf.Pprof.Enabled {
		go func() {
			runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
			runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
			http.ListenAndServe(fmt.Sprintf(":%d", botconf.Pprof.Port), nil)
		}()
	}

	// start
	gamedata.Load()
	job.Load()

	handler.Init()
	app.BotMgr.Start()

	// wait
	<-quit

	// stop
	app.BotMgr.Stop()

	// close log
	log.Notice("bot stopped")
	log.Close()
}
