package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ============================================================================

type Server struct {
	svr           *http.Server        // http server
	mux           *http.ServeMux      // http mux
	upgrader      *websocket.Upgrader // ws upgrader
	behind_proxy  bool                // behind proxy ?
	cb_connection func(sock *Socket)  // callback: connection
	cb_error      func(err error)     // callback: error

	wg sync.WaitGroup
}

// ============================================================================

func CreateServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
		},
	}
}

// ============================================================================

func (self *Server) SetReadBufferSize(n int) *Server {
	self.upgrader.ReadBufferSize = n
	return self
}

func (self *Server) SetWriteBufferSize(n int) *Server {
	self.upgrader.WriteBufferSize = n
	return self
}

func (self *Server) SetSocketUrl(url string) *Server {
	self.mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		self.accept(w, req)
	})

	return self
}

func (self *Server) MapStatic(url string, dir string) *Server {
	self.mux.Handle(url, http.StripPrefix(url, http.FileServer(http.Dir(dir))))
	return self
}

func (self *Server) CheckOrigin(b bool) *Server {
	if b {
		self.upgrader.CheckOrigin = nil
	} else {
		self.upgrader.CheckOrigin = func(req *http.Request) bool {
			return true
		}
	}

	return self
}

func (self *Server) BehindProxy(b bool) *Server {
	self.behind_proxy = b
	return self
}

func (self *Server) Listen(addr string) *Server {
	if self.svr != nil {
		return self
	}

	self.svr = &http.Server{
		Addr:    addr,
		Handler: self.mux,
	}

	self.wg.Add(1)
	go func() {
		defer func() {
			self.svr = nil
			self.wg.Done()
		}()

		err := self.svr.ListenAndServe()
		if err != http.ErrServerClosed {
			if self.cb_error != nil {
				self.cb_error(err)
			}
		}
	}()

	return self
}

func (self *Server) Stop() {
	if self.svr == nil {
		return
	}

	self.svr.Close()
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

func (self *Server) accept(w http.ResponseWriter, req *http.Request) {
	// upgrade
	c, err := self.upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	// create socket
	sock := new_socket(c)

	// parse remote addr
	var fip string
	if self.behind_proxy {
		fip = req.Header.Get("X-Forwarded-For")
	}
	sock.parse_remote_addr(fip)

	// event: connection
	if self.cb_connection != nil {
		self.cb_connection(sock)
	}

	// go rw threads
	go sock.thr_read()
	go sock.thr_write()
}
