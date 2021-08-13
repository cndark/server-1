package tcp

import (
	"sync"
	"time"
)

// ============================================================================

type ConnectQ struct {
	q    chan *connect_f
	wg   sync.WaitGroup
	quit chan int
}

type connect_f struct {
	f     func(done func()) // connect function
	delay int32             // delay in ms
}

// ============================================================================

func NewConnectQ() *ConnectQ {
	return &ConnectQ{
		q:    make(chan *connect_f, 65535),
		quit: make(chan int),
	}
}

func (self *ConnectQ) Open() {
	self.wg.Add(1)
	go self.thr_connect()
}

func (self *ConnectQ) Close() {
	close(self.q)
	close(self.quit)
	self.wg.Wait()
}

func (self *ConnectQ) Connect(f func(done func()), delay int32) {
	// ignore EPIPE
	defer func() { recover() }()

	self.q <- &connect_f{f, delay}
}

// ============================================================================

func (self *ConnectQ) thr_connect() {
	defer self.wg.Done()

	for cf := range self.q {
		cf := cf

		self.wg.Add(1)
		go func() {
			select {
			case <-self.quit:
				self.wg.Done()
				return

			case <-time.After(time.Duration(cf.delay) * time.Millisecond):
				cf.f(self.wg.Done)
			}
		}()
	}
}
