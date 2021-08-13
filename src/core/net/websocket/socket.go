package websocket

import (
	"fmt"
	"fw/src/core"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================================================

type Socket struct {
	c         *websocket.Conn  // underlying connection
	qw        chan *[]byte     // write queue
	rip       string           // remote ip
	rport     int32            // remote port
	wbuf_size int              // write buffer size
	hb        uint32           // heart beat in milliseconds
	cb_data   func(buf []byte) // callback: data
	cb_close  func()           // callback: close
}

// ============================================================================

func new_socket(c *websocket.Conn) *Socket {
	return &Socket{
		c:         c,
		qw:        make(chan *[]byte, 10000),
		wbuf_size: 4096,
	}
}

// ============================================================================

func (self *Socket) Close() {
	self.c.Close()
}

func (self *Socket) Send(buf []byte) {
	// ignore EPIPE
	defer func() { recover() }()

	// push to qw. kick if qw if full
	select {
	case self.qw <- &buf:
	default:
		self.Close()
	}
}

func (self *Socket) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", self.rip, self.rport)
}

func (self *Socket) RemoteIP() string {
	return self.rip
}

func (self *Socket) RemotePort() int32 {
	return self.rport
}

func (self *Socket) SetWriteQSize(n int) {
	close(self.qw)
	self.qw = make(chan *[]byte, n)
}

func (self *Socket) SetWriteBufferSize(n int) {
	self.wbuf_size = n
}

func (self *Socket) HeartBeat(ms uint32) {
	self.hb = ms
}

func (self *Socket) OnData(f func(buf []byte)) {
	self.cb_data = f
}

func (self *Socket) OnClose(f func()) {
	self.cb_close = f
}

// ============================================================================

func (self *Socket) parse_remote_addr(forward_ip string) {
	if forward_ip == "" {
		addr := self.c.RemoteAddr().String()
		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			self.rip = host
			self.rport = core.Atoi32(port)
		}
	} else {
		self.rip = forward_ip
		self.rport = 0
	}
}

func (self *Socket) thr_read() {
	defer func() {
		self.Close()
		close(self.qw)
		self.cb_data = nil
		self.cb_close = nil
	}()

	c := self.c

	for {
		// set deadline
		if self.hb > 0 {
			c.SetReadDeadline(time.Now().Add(time.Duration(self.hb) * time.Millisecond))
		}

		// read
		_, buf, err := c.ReadMessage()
		if err != nil {
			// event: close
			if self.cb_close != nil {
				self.cb_close()
			}

			break
		}

		// event: data
		if self.cb_data != nil {
			self.cb_data(buf)
		}
	}
}

func (self *Socket) thr_write() {
	c := self.c
	qw := self.qw
	// buf := make([]byte, 0, self.wbuf_size)

	for {
		select {
		case rec, ok := <-qw:
			if !ok {
				return
			}

			// buf = append(buf[:0], *rec...)
			// L := len(qw)
			// for L > 0 && len(buf) < self.wbuf_size {
			// 	buf = append(buf, *<-qw...)
			// 	L--
			// }

			err := c.WriteMessage(websocket.BinaryMessage, *rec)
			if err != nil {
				return
			}
		}
	}
}
