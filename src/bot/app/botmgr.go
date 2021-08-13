package app

import (
	"fw/src/bot/botconf"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/net/websocket"
	"fw/src/core/packet"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================

var BotMgr = &botmgr_t{
	bots: make(map[int32]*Bot),
}

// ============================================================================

type botmgr_t struct {
	bots    map[int32]*Bot
	bot_cnt int32 // bot count (atomic counter)

	locker sync.Mutex
	wg     sync.WaitGroup
}

// ============================================================================

func (self *botmgr_t) Start() {
	i := 0

	for _, e := range botconf.Bots {
		e := e

		// opt
		opt := &bot_opt{
			svr:        e.Svr,
			sdk:        e.Sdk,
			model:      e.Model,
			job_prefix: e.JobPrefix,
			job_itv:    e.JobItv,
			args:       parse_args(e.Args),
			grp:        e.Grp,
		}

		// do connect
		for id := e.From; id <= e.To; id++ {
			id := id

			self.wg.Add(2)
			go func() {
				defer self.wg.Done()

				i++
				i := i

				time.Sleep(time.Duration(i*50) * time.Millisecond)

				if strings.HasPrefix(e.Addr, "ws://") || strings.HasPrefix(e.Addr, "wss://") {
					// websocket
					websocket.Connect(e.Addr, 3000, func(err error, sock *websocket.Socket) {
						defer self.wg.Done()
						self.on_connect(err, sock, id, opt)
					})
				} else {
					// tcp
					tcp.Connect(e.Addr, 3000, func(err error, sock *tcp.Socket) {
						defer self.wg.Done()
						self.on_connect(err, sock, id, opt)
					})
				}
			}()
		}
	}
}

func (self *botmgr_t) Stop() {
	// wait for all connect i/o
	self.wg.Wait()

	// close all bots
	self.close_all_bots()

	// wait
	for {
		if self.bot_count() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (self *botmgr_t) on_connect(err error, sock socket_t, id int32, opt *bot_opt) {
	if err != nil {
		log.Error("connect to server failed:", err)
		return
	}

	// add session
	bot := new_bot(sock, id, opt)
	self.add_bot(bot)

	// sock event: data
	sock.OnData(func(buf []byte) {
		var p packet.Packet
		var err error

		for len(buf) > 0 {
			// read packet
			p, buf, err = bot.preader.Read(buf)
			if err != nil {
				log.Warning("reading packet failed:", err)
				sock.Close()
				return
			}

			// no packet yet
			if p == nil {
				return
			}

			// got packet. dispatch
			bot.Dispatch(p)
		}
	})

	// sock event: close
	sock.OnClose(func() {
		// remove bot
		self.remove_bot(bot)
	})

	// run
	bot.run()
}

func (self *botmgr_t) bot_count() int32 {
	return atomic.LoadInt32(&self.bot_cnt)
}

// ============================================================================

func (self *botmgr_t) add_bot(bot *Bot) {
	self.locker.Lock()
	self.bots[bot.Id] = bot
	atomic.AddInt32(&self.bot_cnt, 1)
	self.locker.Unlock()

	// bot event: connected
	bot.Push(func() {
		bot.OnConnected()
	})
}

func (self *botmgr_t) remove_bot(bot *Bot) {
	nobot := false

	self.locker.Lock()
	delete(self.bots, bot.Id)
	atomic.AddInt32(&self.bot_cnt, -1)
	if self.bot_cnt == 0 {
		nobot = true
	}
	self.locker.Unlock()

	// bot event: disconnected
	bot.Push(func() {
		bot.OnDisconnected()

		// stop after disconnected event is consumed
		bot.stop()

		// end-program check after disconnected event is consumed
		if nobot {
			evtmgr.Fire("all.bot.dies")
		}
	})
}

func (self *botmgr_t) close_all_bots() {
	self.locker.Lock()
	defer self.locker.Unlock()

	for _, bot := range self.bots {
		bot.Close()
	}
}
