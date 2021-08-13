package app

import (
	"fw/src/core"
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/packet"
	"fw/src/core/sched/loop"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"sort"
	"sync"
	"time"
)

// ============================================================================

const (
	C_max_gateid = 100000
)

// ============================================================================

var seq_gateid int32 = C_max_gateid
var use_router_id int32

// ============================================================================

var NetMgr = &netmgr_t{
	gates:    make(map[int32]*SocketGW),
	connectq: tcp.NewConnectQ(),
	cnn_rt:   make(map[int32]*SocketRt),
}

// ============================================================================

type netmgr_t struct {
	svr4gw *tcp.Server         // server for gate
	gates  map[int32]*SocketGW // gate map

	connectq *tcp.ConnectQ       // connect queue
	cnn_rt   map[int32]*SocketRt // router map

	locker_gw sync.Mutex
	locker_rt sync.Mutex
}

// ============================================================================

func (self *netmgr_t) Start() {
	log.Info("starting net mgr ...")

	self.listen_on_gates()

	self.connectq.Open()
	self.connect_to_routers()

	use_router_id = self.choose_router()
}

func (self *netmgr_t) Stop() {
	log.Info("stopping net mgr ...")

	// stop servers
	self.svr4gw.Stop()
	self.close_all_gates()

	// stop connections
	self.connectq.Close()
	self.close_all_connections()

	// wait
	for {
		if self.count_gw() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (self *netmgr_t) RegisterGate(gw *SocketGW, id int32) bool {
	self.locker_gw.Lock()
	defer self.locker_gw.Unlock()

	// check reg id
	if id >= C_max_gateid {
		return false
	}

	// old entry MUST exist
	// new entry MUST NOT exist
	if self.gates[gw.id] == nil || self.gates[id] != nil {
		return false
	}

	// remove old entry
	delete(self.gates, gw.id)

	// set registered id
	gw.id = id

	// add new entry
	self.gates[gw.id] = gw

	log.Notice("gate registered:", id)

	return true
}

func (self *netmgr_t) Send2Gate(gateid int32, message msg.Message) {
	if gateid == 0 {
		// broadcast
		for _, gw := range self.reg_gates_array() {
			gw.SendMsg(message)
		}
	} else {
		// specific one
		gw := self.find_gate(gateid)
		if gw != nil {
			gw.SendMsg(message)
		}
	}
}

func (self *netmgr_t) Send2Player(sid uint64, message msg.Message) {
	// find gate
	gateid := sid_to_gateid(sid)
	gw := self.find_gate(gateid)
	if gw == nil {
		return
	}

	// marshal
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	// assemble
	p := packet.Assemble_B8(message.MsgId(), body, sid)

	// send
	gw.SendPacket(p)
}

func (self *netmgr_t) Forward2Player(sid uint64, p packet.Packet) {
	// find gate
	gateid := sid_to_gateid(sid)
	gw := self.find_gate(gateid)
	if gw == nil {
		return
	}

	// add sid
	p.Add_B8(sid)

	// send
	gw.SendPacket(p)
}

func (self *netmgr_t) Send2Game(svrid int32, message msg.Message) {
	rt := self.use_router()
	if rt != nil {
		p := rt.to_packet(message, uint64(1)<<32|uint64(svrid), "")
		if p != nil {
			rt.SendPacket(p)
		}
	}
}

func (self *netmgr_t) Send2CrossPlayer(svrid int32, plrid string, message msg.Message) {
	rt := self.use_router()
	if rt != nil {
		p := rt.to_packet(message, uint64(1)<<32|uint64(svrid), plrid)
		if p != nil {
			rt.SendPacket(p)
		}
	}
}

func (self *netmgr_t) check_route(p packet.Packet) packet.Packet {
	_, plrid := p.Remove_B8_Str()
	if plrid == "" {
		return p
	}

	loop.Push(func() {
		plr := PlayerMgr.FindPlayerById(plrid)
		if plr != nil && plr.IsOnline() {
			NetMgr.Forward2Player(plr.sid, p)
		}
	})

	return nil
}

func (self *netmgr_t) count_gw() int {
	self.locker_gw.Lock()
	defer self.locker_gw.Unlock()

	return len(self.gates)
}

// ============================================================================

func (self *netmgr_t) listen_on_gates() {
	self.svr4gw = tcp.CreateServer().
		OnConnection(func(sock *tcp.Socket) {
			// set sock opt
			sock.SetReadBufferSize(512000)
			sock.SetWriteBufferSize(512000)
			sock.SetWriteQSize(3_000_000)

			// add gate
			gw := new_socket_gw(sock)
			self.add_gate(gw)

			// sock event: data
			sock.OnData(func(buf []byte) {
				var p packet.Packet
				var err error

				for len(buf) > 0 {
					// read packet
					p, buf, err = gw.preader.Read(buf)
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
					gw.Dispatch(p)
				}
			})

			// sock event: close
			sock.OnClose(func() {
				// remove gate
				self.remove_gate(gw)
			})
		}).
		OnError(func(err error) {
			core.Panic("listen on gate failed:", err)
		}).
		Listen(config.CurGame.Addr4GW)

	log.Notice("listen on gates:", config.CurGame.Addr4GW)
}

func (self *netmgr_t) add_gate(gw *SocketGW) {
	self.locker_gw.Lock()
	defer self.locker_gw.Unlock()

	self.gates[gw.id] = gw
}

func (self *netmgr_t) remove_gate(gw *SocketGW) {
	if gw.is_registered() {
		log.Warning("gate dropped:", gw.id)

		loop.Push(func() {
			// offline players from this gate
			PlayerMgr.OfflinePlayers(gw.id)
		})
	}

	self.locker_gw.Lock()
	delete(self.gates, gw.id)
	self.locker_gw.Unlock()
}

func (self *netmgr_t) find_gate(id int32) *SocketGW {
	self.locker_gw.Lock()
	defer self.locker_gw.Unlock()

	gw := self.gates[id]
	if gw != nil && gw.is_registered() {
		return gw
	} else {
		return nil
	}
}

func (self *netmgr_t) reg_gates_array() (ret []*SocketGW) {
	self.locker_gw.Lock()
	defer self.locker_gw.Unlock()

	for _, gw := range self.gates {
		if gw.is_registered() {
			ret = append(ret, gw)
		}
	}

	return
}

func (self *netmgr_t) close_all_gates() {
	var arr []*SocketGW

	self.locker_gw.Lock()
	for _, cnn := range self.gates {
		arr = append(arr, cnn)
	}
	self.locker_gw.Unlock()

	for _, v := range arr {
		v.Close()
	}
}

// ============================================================================

func (self *netmgr_t) connect_to_routers() {
	log.Info("connecting to routers ...")

	m := config.Routers
	for _, conf := range m {
		conf := conf

		var f_cnn func(done func())

		f_cnn = func(done func()) {
			tcp.Connect(conf.Addr4C, 3000, func(err error, sock *tcp.Socket) {
				defer done()

				if err != nil {
					log.Warning("connect to router failed:", conf.Name, err)
					self.connectq.Connect(f_cnn, 5000)
					return
				}

				// set sock opt
				sock.SetReadBufferSize(512000)
				sock.SetWriteBufferSize(512000)
				sock.SetWriteQSize(3_000_000)

				// create cnn
				cnn_rt := new_socket_rt(sock, conf.Id)
				self.add_router(cnn_rt)

				// sock event: data
				sock.OnData(func(buf []byte) {
					var p packet.Packet
					var err error

					for len(buf) > 0 {
						// read packet
						p, buf, err = cnn_rt.preader.Read(buf)
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
						cnn_rt.Dispatch(p)
					}
				})

				// sock event: close
				sock.OnClose(func() {
					log.Warning("connection to router disconnected:", conf.Name)

					self.remove_router(cnn_rt)

					self.connectq.Connect(f_cnn, 5000)
					return
				})

				// open
				cnn_rt.Open()
			})
		}

		self.connectq.Connect(f_cnn, 0)
	}
}

func (self *netmgr_t) add_router(rt *SocketRt) {
	self.locker_rt.Lock()
	defer self.locker_rt.Unlock()

	self.cnn_rt[rt.Id] = rt
}

func (self *netmgr_t) remove_router(rt *SocketRt) {
	self.locker_rt.Lock()
	defer self.locker_rt.Unlock()

	delete(self.cnn_rt, rt.Id)
}

func (self *netmgr_t) use_router() *SocketRt {
	self.locker_rt.Lock()
	defer self.locker_rt.Unlock()

	return self.cnn_rt[use_router_id]
}

func (self *netmgr_t) close_all_connections() {
	// close cnn_rt
	var arr []*SocketRt

	self.locker_rt.Lock()
	for _, rt := range self.cnn_rt {
		arr = append(arr, rt)
	}
	self.locker_rt.Unlock()

	for _, v := range arr {
		v.Close()
	}
}

func (self *netmgr_t) choose_router() int32 {
	var arr []int32

	m := config.Routers // in case of config update in another thread
	for _, v := range m {
		arr = append(arr, v.Id)
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	return arr[config.CurGame.Id%int32(len(arr))]
}
