package app

import (
	"fw/src/core"
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/packet"
	"fw/src/shared/config"
	"sync"
	"time"
)

// ============================================================================

const (
	C_max_cid = 100000
)

// ============================================================================

var seq_cid int32 = C_max_cid

// ============================================================================

var NetMgr = &netmgr_t{
	map_c: make(map[int32]*SocketC),
}

// ============================================================================

type netmgr_t struct {
	svr4c *tcp.Server        // server for c
	map_c map[int32]*SocketC // c map

	locker sync.Mutex
}

// ============================================================================

func (self *netmgr_t) Start() {
	log.Info("starting net mgr ...")

	self.listen_on_c()
}

func (self *netmgr_t) Stop() {
	log.Info("stopping net mgr ...")

	// stop servers
	self.svr4c.Stop()
	self.close_all_c()

	// wait
	for {
		if self.count_c() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (self *netmgr_t) RegisterC(c *SocketC, id int32) bool {
	self.locker.Lock()
	defer self.locker.Unlock()

	// check reg id
	if id >= C_max_cid {
		return false
	}

	// old entry MUST exist
	// new entry MUST NOT exist
	if self.map_c[c.id] == nil || self.map_c[id] != nil {
		return false
	}

	// remove old entry
	delete(self.map_c, c.id)

	// set registered id
	c.id = id

	// add new entry
	self.map_c[c.id] = c

	log.Notice("c registered:", id)

	return true
}

func (self *netmgr_t) check_route(p packet.Packet) packet.Packet {
	b8, _ := p.Peek_B8_Str()
	if b8 == 0 {
		p.Remove_B8_Str()
		return p
	}

	who := b8 >> 32
	id := int32(b8 & 0x00000000ffffffff)

	switch who {
	case 1: // game
		c := NetMgr.find_c(id)
		if c != nil {
			c.SendPacket(p)
		}
	}

	return nil
}

func (self *netmgr_t) count_c() int {
	self.locker.Lock()
	defer self.locker.Unlock()

	return len(self.map_c)
}

// ============================================================================

func (self *netmgr_t) listen_on_c() {
	self.svr4c = tcp.CreateServer().
		OnConnection(func(sock *tcp.Socket) {
			// set sock opt
			sock.SetReadBufferSize(512000)
			sock.SetWriteBufferSize(512000)
			sock.SetWriteQSize(3_000_000)

			// add c
			c := new_socket_c(sock)
			self.add_c(c)

			// sock event: data
			sock.OnData(func(buf []byte) {
				var p packet.Packet
				var err error

				for len(buf) > 0 {
					// read packet
					p, buf, err = c.preader.Read(buf)
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
					c.Dispatch(p)
				}
			})

			// sock event: close
			sock.OnClose(func() {
				// remove c
				self.remove_c(c)
			})
		}).
		OnError(func(err error) {
			core.Panic("listen on c failed:", err)
		}).
		Listen(config.CurRouter.Addr4C)

	log.Notice("listen on c:", config.CurRouter.Addr4C)
}

func (self *netmgr_t) add_c(c *SocketC) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.map_c[c.id] = c
}

func (self *netmgr_t) remove_c(c *SocketC) {
	self.locker.Lock()
	defer self.locker.Unlock()

	delete(self.map_c, c.id)
}

func (self *netmgr_t) find_c(id int32) *SocketC {
	self.locker.Lock()
	defer self.locker.Unlock()

	c := self.map_c[id]
	if c != nil && c.is_registered() {
		return c
	} else {
		return nil
	}
}

func (self *netmgr_t) close_all_c() {
	var arr []*SocketC

	self.locker.Lock()
	for _, c := range self.map_c {
		arr = append(arr, c)
	}
	self.locker.Unlock()

	for _, v := range arr {
		v.Close()
	}
}
