package tcp

import (
	"net"
	"time"
)

// ============================================================================

func Connect(addr string, timeout int32, f func(err error, sock *Socket)) {
	go func() {
		c, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			if f != nil {
				f(err, nil)
			}
			return
		}

		// create socket
		sock := new_socket(c)

		// parse remote addr
		sock.parse_remote_addr()

		// event: connect
		if f != nil {
			f(nil, sock)
		}

		// go rw threads
		go sock.thr_read()
		go sock.thr_write()
	}()
}
