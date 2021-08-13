package app

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/net/websocket"
	"fw/src/core/packet"
	"fw/src/gate/msg"
	"fw/src/shared/config"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================

var NetMgr = &netmgr_t{
	sessions:       make(map[uint64]*Session),
	connectq:       tcp.NewConnectQ(),
	cnn_gs:         make(map[int32]*SocketGS),
	sched_gs_names: make(map[string]bool),
}

// ============================================================================

const (
	C_max_session_count = 30000
)

// ============================================================================

type netmgr_t struct {
	svr4c    *tcp.Server         // server for client
	ws4c     *websocket.Server   // websocket server for client
	sessions map[uint64]*Session // session map
	sess_cnt int32               // session count (atomic counter)

	connectq *tcp.ConnectQ       // connect queue
	cnn_gs   map[int32]*SocketGS // gs map

	locker_c  sync.Mutex
	locker_gs sync.Mutex

	// current gs ids that are scheduled for connection
	sched_gs_names map[string]bool
	locker_sgn     sync.Mutex
}

// ============================================================================

func (self *netmgr_t) Start() {
	log.Info("starting net mgr ...")

	self.listen_on_client()
	self.listen_on_client_ws()

	self.connectq.Open()
	self.connect_to_games()
}

func (self *netmgr_t) Stop() {
	log.Info("stopping net mgr ...")

	// stop servers
	self.svr4c.Stop()
	self.ws4c.Stop()
	self.close_sessions(0)

	// stop connections
	self.connectq.Close()
	self.close_all_connections()

	// wait
	for {
		if self.session_count() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (self *netmgr_t) Forward2GS(gsid int32, sid uint64, p packet.Packet) {
	gs := self.find_game(gsid)
	if gs != nil {
		p.Add_B8(sid)
		gs.SendPacket(p)
	}
}

func (self *netmgr_t) Forward2Session(sid uint64, p packet.Packet) {
	sess := self.FindSession(sid)
	if sess != nil {
		sess.SendPacket(p)
	}
}

func (self *netmgr_t) Send2GS(gsid int32, message msg.Message) bool {
	gs := self.find_game(gsid)
	if gs != nil {
		gs.SendMsg(message)
		return true
	}

	return false
}

func (self *netmgr_t) ModifyGameConnections() {
	self.connect_to_games()
}

func (self *netmgr_t) KickGamesPlayer() {
	fileName := "./KICK_GAMEIDS"

	defer os.Remove(fileName)
	d, err := ioutil.ReadFile(fileName)
	if err != nil || len(d) == 0 {
		return
	}

	for _, v := range strings.Split(string(d), " ") {
		gsid := core.Atoi32(strings.Trim(v, " \n"))
		if gsid > 0 && gsid <= config.GameIdMax {
			self.close_sessions(gsid)
		}
	}
}

// ============================================================================

func (self *netmgr_t) listen_on_client() {
	// some cloud-servers DO NOT allow us to listen on specific WAN IP
	// 	so, listen on 0.0.0.0
	addr := fmt.Sprintf(":%d", config.CurGate.Port)

	self.svr4c = tcp.CreateServer().
		OnConnection(func(sock *tcp.Socket) {
			// check session count
			if self.session_count() > C_max_session_count {
				sock.Close()
				return
			}

			// 1min heartbeat kick
			sock.HeartBeat(60000)

			// add session
			sess := new_session(sock)
			self.add_session(sess)

			// sock event: data
			sock.OnData(func(buf []byte) {
				var p packet.Packet
				var err error

				for len(buf) > 0 {
					// read packet
					p, buf, err = sess.preader.Read(buf)
					if err != nil {
						log.Debug("reading packet failed:", sock.RemoteAddr(), err)
						sock.Close()
						return
					}

					// no packet yet
					if p == nil {
						return
					}

					// got packet. dispatch
					sess.Dispatch(p)
				}
			})

			// sock event: close
			sock.OnClose(func() {
				// remove session
				self.remove_session(sess)
			})
		}).
		OnError(func(err error) {
			core.Panic("listen on client failed:", err)
		}).
		Listen(addr)

	log.Notice("listen on client:", addr)
}

func (self *netmgr_t) listen_on_client_ws() {
	// some cloud-servers DO NOT allow us to listen on specific WAN IP
	// 	so, listen on 0.0.0.0
	addr := fmt.Sprintf(":%d", config.CurGate.WsPort)

	self.ws4c = websocket.CreateServer().
		SetSocketUrl("/gate").
		MapStatic("/h5/", "./h5").
		CheckOrigin(false).
		BehindProxy(config.Common.BehindProxy).
		OnConnection(func(sock *websocket.Socket) {
			// check session count
			if self.session_count() > C_max_session_count {
				sock.Close()
				return
			}

			// 1min heartbeat kick
			sock.HeartBeat(60000)

			// add session
			sess := new_session(sock)
			self.add_session(sess)

			// sock event: data
			sock.OnData(func(buf []byte) {
				var p packet.Packet
				var err error

				for len(buf) > 0 {
					// read packet
					p, buf, err = sess.preader.Read(buf)
					if err != nil {
						log.Debug("reading packet failed:", sock.RemoteAddr(), err)
						sock.Close()
						return
					}

					// no packet yet
					if p == nil {
						return
					}

					// got packet. dispatch
					sess.Dispatch(p)
				}
			})

			// sock event: close
			sock.OnClose(func() {
				// remove session
				self.remove_session(sess)
			})
		}).
		OnError(func(err error) {
			core.Panic("listen on client failed:", err)
		}).
		Listen(addr)

	log.Notice("listen on client for websocket:", addr)
}

func (self *netmgr_t) add_session(sess *Session) {
	self.locker_c.Lock()
	defer self.locker_c.Unlock()

	self.sessions[sess.id] = sess
	atomic.AddInt32(&self.sess_cnt, 1)
}

func (self *netmgr_t) remove_session(sess *Session) {
	self.locker_c.Lock()
	delete(self.sessions, sess.id)
	atomic.AddInt32(&self.sess_cnt, -1)
	self.locker_c.Unlock()

	// logout player
	sess.LogoutPlayer()
}

func (self *netmgr_t) FindSession(sid uint64) *Session {
	self.locker_c.Lock()
	defer self.locker_c.Unlock()

	return self.sessions[sid]
}

func (self *netmgr_t) close_sessions(gsid int32) {
	var arr []*Session

	self.locker_c.Lock()
	if gsid == 0 {
		// all sessions
		for _, sess := range self.sessions {
			arr = append(arr, sess)
		}
	} else {
		// sessions to gsid
		for _, sess := range self.sessions {
			if sess.gsid == gsid {
				arr = append(arr, sess)
			}
		}
	}
	self.locker_c.Unlock()

	for _, v := range arr {
		v.Close()
	}
}

func (self *netmgr_t) session_count() int32 {
	return atomic.LoadInt32(&self.sess_cnt)
}

// ============================================================================

func (self *netmgr_t) connect_to_games() {
	log.Info("connecting to games ...")

	self.locker_sgn.Lock()
	defer self.locker_sgn.Unlock()

	// get current games
	m := config.Games

	// unsched removed games
	{
		arr := make([]string, 0, len(self.sched_gs_names))
		for name, _ := range self.sched_gs_names {
			if m[name] == nil {
				arr = append(arr, name)
			}
		}

		for _, name := range arr {
			delete(self.sched_gs_names, name)
		}
	}

	// sched new games
	for _, conf := range m {
		conf := conf

		// check if already scheduled
		if self.sched_gs_names[conf.Name] {
			continue
		} else {
			self.sched_gs_names[conf.Name] = true
		}

		// connect function
		var f_cnn func(done func())

		f_cnn = func(done func()) {
			tcp.Connect(conf.Addr4GW, 3000, func(err error, sock *tcp.Socket) {
				defer done()

				if err != nil {
					log.Warning("connect to game failed:", conf.Name, err)

					// reconnect only if the game is NOT removed
					self.locker_sgn.Lock()
					b := self.sched_gs_names[conf.Name]
					self.locker_sgn.Unlock()

					if b {
						self.connectq.Connect(f_cnn, 5000)
					}

					return
				}

				// set sock opt
				sock.SetReadBufferSize(512000)
				sock.SetWriteBufferSize(512000)
				sock.SetWriteQSize(3_000_000)

				// create cnn
				cnn_gs := new_socket_gs(sock, conf.Id)
				self.add_game(cnn_gs)

				// sock event: data
				sock.OnData(func(buf []byte) {
					var p packet.Packet
					var err error

					for len(buf) > 0 {
						// read packet
						p, buf, err = cnn_gs.preader.Read(buf)
						if err != nil {
							log.Debug("reading packet failed:", sock.RemoteAddr(), err)
							sock.Close()
							return
						}

						// no packet yet
						if p == nil {
							return
						}

						// got packet. dispatch
						cnn_gs.Dispatch(p)
					}
				})

				// sock event: close
				sock.OnClose(func() {
					log.Warning("connection to gs disconnected:", conf.Name)

					self.remove_game(cnn_gs)

					self.connectq.Connect(f_cnn, 5000)
					return
				})

				// open
				cnn_gs.Open()
			})
		}

		// connect
		self.connectq.Connect(f_cnn, 0)
	}
}

func (self *netmgr_t) add_game(gs *SocketGS) {
	self.locker_gs.Lock()
	defer self.locker_gs.Unlock()

	self.cnn_gs[gs.Id] = gs
}

func (self *netmgr_t) remove_game(gs *SocketGS) {
	self.locker_gs.Lock()
	delete(self.cnn_gs, gs.Id)
	defer self.locker_gs.Unlock()

	self.close_sessions(gs.Id)
}

func (self *netmgr_t) find_game(gsid int32) *SocketGS {
	self.locker_gs.Lock()
	defer self.locker_gs.Unlock()

	return self.cnn_gs[gsid]
}

func (self *netmgr_t) close_all_connections() {
	// close cnn_gs
	var arr []*SocketGS

	self.locker_gs.Lock()
	for _, gs := range self.cnn_gs {
		arr = append(arr, gs)
	}
	self.locker_gs.Unlock()

	for _, v := range arr {
		v.Close()
	}
}
