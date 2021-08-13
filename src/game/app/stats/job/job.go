package job

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/sched/loop"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"time"
)

// ============================================================================

type online_t struct {
	avg  float64
	peek int32
	n    int32

	last_t time.Time
}

// ============================================================================

func init() {
	(&online_t{
		last_t: core.StartOfDay(time.Now()),
	}).start()
}

// ============================================================================

func (self *online_t) start() {
	loop.SetTimeout(time.Now().Add(time.Minute*1), func() {
		self.chk_rst()
		self.stats()
		self.start()
	})
}

func (self *online_t) chk_rst() {
	// check
	t := core.StartOfDay(time.Now())
	if t.Equal(self.last_t) {
		return
	}

	// reset
	self.last_t = t

	self.peek = 0
	self.avg = 0
	self.n = 0
}

func (self *online_t) stats() {
	v := app.PlayerMgr.NumOnline()

	// avg
	self.avg = (float64(self.n)*self.avg + float64(v)) / float64(self.n+1)
	self.n++

	// peek
	if v > self.peek {
		self.peek = v
	}

	// fire
	evtmgr.Fire(gconst.Evt_OnlineNum, v, int32(self.avg), self.peek)
}
