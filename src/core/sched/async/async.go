package async

import (
	"fw/src/core"
	"fw/src/core/log"
	"sync"
)

// ============================================================================

var (
	qn   int
	qs   []*async_queue_t
	quit = make(chan int)
	wg   sync.WaitGroup
)

// ============================================================================

type async_queue_t struct {
	ch chan func() // channel
	w  int         // workers
}

// ============================================================================

// ws: worker count of each queue
func Init(ws []int) {
	qn = len(ws)
	qs = make([]*async_queue_t, 0, qn)

	for _, n := range ws {
		qs = append(qs, &async_queue_t{
			ch: make(chan func(), 1_000_000),
			w:  n,
		})
	}
}

func Start() {
	// each queue
	for _, q := range qs {
		q := q

		// spanw workers
		for i := 0; i < q.w; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					select {
					case <-quit:
						return

					case f := <-q.ch:
						core.PCall(f)
					}
				}
			}()
		}
	}
}

func Stop() {
	close(quit)
	wg.Wait()

	for _, q := range qs {
		close(q.ch)
		for f := range q.ch {
			f()
		}
	}
}

// default push -> queue 0
func Push(f func()) {
	// ignore EPIPE
	defer func() { recover() }()

	select {
	case qs[0].ch <- f:
	default:
		log.Error("async queue [0] is FULL. push discarded")
	}
}

// specific push -> queue 'n'
func PushQ(n int, f func()) {
	// ignore EPIPE
	defer func() { recover() }()

	if n < 0 || n >= qn {
		return
	}

	select {
	case qs[n].ch <- f:
	default:
		log.Error("async queue [n] is FULL. push discarded", n)
	}
}

func QLen(n int) int32 {
	if n < 0 || n >= qn {
		return 0
	}

	return int32(len(qs[n].ch))
}
