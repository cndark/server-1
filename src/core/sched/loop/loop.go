package loop

import (
	"fw/src/core"
	"fw/src/core/log"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================

var (
	q      = make(chan func(), 1000000)
	timerq = core.NewTimerQueue()
	quit   = make(chan int)
	wg     sync.WaitGroup

	num_handled int32
)

// ============================================================================

func Run() {
	wg.Add(1)

	go loop_f()
	go loop_timer()
}

func Stop() {
	close(quit)
	wg.Wait()
}

func Push(f func()) {
	// ignore EPIPE
	defer func() { recover() }()

	select {
	case q <- f:
	default:
		log.Error("loopq FULL. push discarded")
	}
}

func SetTimeout(ts time.Time, f func()) *core.Timer {
	return timerq.SetTimeout(ts, f)
}

func CancelTimer(t *core.Timer) {
	timerq.Cancel(t)
}

func UpdateTimer(t *core.Timer, ts time.Time) {
	timerq.Update(t, ts)
}

func QLen() int32 {
	return int32(len(q))
}

func NumHandled() int32 {
	return atomic.SwapInt32(&num_handled, 0)
}

// ============================================================================

func loop_f() {
	defer wg.Done()

	for f := range q {
		core.PCall(f)
		atomic.AddInt32(&num_handled, 1)
	}
}

func loop_timer() {
	defer close(q)

	for {
		select {
		case <-quit:
			return

		default:
			Push(func() {
				now := time.Now()
				for timerq.Expire(now) {
				}
			})

			time.Sleep(200 * time.Millisecond)
		}
	}
}

// ============================================================================

func TimeSlice(arr_len int, slice_n int, delay_ms int32, f func(int), done ...func()) {
	now := time.Now()
	n := 0

	for i := 0; i < arr_len; i += slice_n {
		i := i
		SetTimeout(now.Add(time.Millisecond*time.Duration(delay_ms)*time.Duration(n)), func() {
			j := i + slice_n
			if j > arr_len {
				j = arr_len
			}

			for k := i; k < j; k++ {
				f(k)
			}
		})
		n++
	}

	if len(done) > 0 {
		SetTimeout(now.Add(time.Millisecond*time.Duration(delay_ms)*time.Duration(n)), func() {
			for _, f2 := range done {
				f2()
			}
		})
	}
}
