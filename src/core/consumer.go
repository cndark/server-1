package core

import "sync"

// ============================================================================

type Consumer struct {
	n  int
	ch chan func()
	wg sync.WaitGroup
}

// ============================================================================

func NewConsumer(n int) *Consumer {
	c := &Consumer{
		n:  n,
		ch: make(chan func(), n),
	}

	c.wg.Add(n)
	for i := 0; i < n; i++ {
		Go(func() {
			defer c.wg.Done()

			for f := range c.ch {
				f()
			}
		})
	}

	return c
}

// ============================================================================

func (self *Consumer) Push(f func()) {
	self.ch <- f
}

func (self *Consumer) Close() {
	close(self.ch)
	self.wg.Wait()
}
