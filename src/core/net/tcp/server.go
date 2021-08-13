package tcp

import (
	"net"
	"sync"
)

// ============================================================================

type Server struct {
	lsn           net.Listener       // underlying listener
	cb_connection func(sock *Socket) // callback: connection
	cb_error      func(err error)    // callback: error

	wg sync.WaitGroup
}

// ============================================================================

func CreateServer() *Server {
	return &Server{}
}

// ============================================================================

func (self *Server) Listen(addr string) *Server {
	if self.lsn != nil {
		return self
	}

	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		if self.cb_error != nil {
			self.cb_error(err)
		}
		return self
	}

	self.lsn = lsn

	self.wg.Add(1)
	go self.thr_accept()

	return self
}

func (self *Server) Stop() {
	if self.lsn == nil {
		return
	}

	self.lsn.Close()
	self.wg.Wait()
}

func (self *Server) OnConnection(f func(sock *Socket)) *Server {
	self.cb_connection = f
	return self
}

func (self *Server) OnError(f func(err error)) *Server {
	self.cb_error = f
	return self
}

// ============================================================================

func (self *Server) thr_accept() {
	defer func() {
		self.lsn.Close()
		self.lsn = nil
		self.wg.Done()
	}()

	for {
		// accept
		c, err := self.lsn.Accept()
		if err != nil {
			break
		}

		// create socket
		sock := new_socket(c)

		// parse remote addr
		sock.parse_remote_addr()

		// event: connection
		if self.cb_connection != nil {
			self.cb_connection(sock)
		}

		// go rw threads
		go sock.thr_read()
		go sock.thr_write()
	}
}
