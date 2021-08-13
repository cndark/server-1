package stage

import (
	"fw/src/core"
	"fw/src/core/sched/loop"
	"time"
)

// ============================================================================

type StageLine struct {
	id     int32
	stages []*Stage

	tid *core.Timer
}

type Stage struct {
	ts time.Time
	f  func()
}

// ============================================================================

func NewStageLine(id int32) *StageLine {
	return &StageLine{
		id: id,
	}
}

func (self *StageLine) Start(pending ...func()) {
	L := len(self.stages)
	if L == 0 {
		return
	}

	// find current stage index
	idx := -1
	now := time.Now()

	for i := L - 1; i >= 0; i-- {
		stg := self.stages[i]

		if !stg.ts.After(now) {
			idx = i
			break
		}
	}

	// sched next timer
	self.sched_next_timer(idx + 1)

	// check current stage
	if idx == -1 {
		// execute pending callback
		for _, pf := range pending {
			pf()
		}
	} else {
		// jump to current stage.
		cur_stg := self.stages[idx]

		if cur_stg.f != nil {
			cur_stg.f()
		}
	}
}

func (self *StageLine) Stop() {
	if self.tid != nil {
		loop.CancelTimer(self.tid)
		self.tid = nil
	}
}

func (self *StageLine) Id() int32 {
	return self.id
}

func (self *StageLine) Add(ts time.Time, f func()) {
	self.stages = append(self.stages, &Stage{
		ts: ts,
		f:  f,
	})
}

func (self *StageLine) Get(i int) *Stage {
	if i < 0 || i >= len(self.stages) {
		return nil
	}

	return self.stages[i]
}

func (self *StageLine) sched_next_timer(idx int) {
	// check if the line is ended
	if idx >= len(self.stages) {
		return
	}

	// get next stage
	next_stg := self.stages[idx]

	// start timer
	self.tid = loop.SetTimeout(next_stg.ts, func() {
		// sched next timer first,
		//  so that we're able to call Stop() from stage callback to end the line
		self.sched_next_timer(idx + 1)

		// stage callback
		if next_stg.f != nil {
			next_stg.f()
		}
	})
}

// ============================================================================

func (self *Stage) GetTs() time.Time {
	return self.ts
}
