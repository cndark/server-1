package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================================================

func Connect(addr string, timeout int32, f func(err error, sock *Socket)) {
	go func() {
		dialer := websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: time.Millisecond * time.Duration(timeout),
		}

		c, _, err := dialer.Dial(addr, nil)
		if err != nil {
			if f != nil {
				f(err, nil)
			}
			return
		}

		// create socket
		sock := new_socket(c)

		// parse remote addr
		sock.parse_remote_addr("")

		// event: connect
		if f != nil {
			f(nil, sock)
		}

		// go rw threads
		go sock.thr_read()
		go sock.thr_write()
	}()
}
