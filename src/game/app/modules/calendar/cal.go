package calendar

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/sched/stage"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/worlddata"
	"strings"
	"time"
)

// ============================================================================

type first_time_t struct {
	t0    time.Time // t0
	unit  string    // repeat time unit
	ldays int32     // loop days
	lstop time.Time // loop stop
	abs   bool      // is absolute time line
}

// ============================================================================

func Start() {
	for _, v := range gamedata.ConfCalendar.Items() {
		v := v

		// check
		L := len(v.StageList)
		if L == 0 {
			continue
		}

		// parse first time (mostly repeat time)
		ft := parse_first_time(v.StageList[0].Ts)
		if ft == nil {
			continue
		}

		// parse last time (end time)
		var t_last time.Time

		if L == 1 {
			t_last = ft.t0
		} else {
			t_last = core.ParseRelativeTime(ft.t0, v.StageList[L-1].Ts)
		}

		// check if we should sched old stage-line or new stage-line
		now := time.Now()
		if !t_last.After(now) {
			if !ft.inc() {
				continue
			}
		}

		// sched function
		var sched_next func(*first_time_t)

		sched_next = func(ft *first_time_t) {
			// parse stage times
			ts := make([]time.Time, 0, L)
			ts = append(ts, ft.t0)

			for i := 1; i < L; i++ {
				ts = append(ts, core.ParseRelativeTime(ft.t0, v.StageList[i].Ts))
			}

			// new stage line
			sl := stage.NewStageLine(0)

			// add stages
			// 	t0: stage-line start time
			// 	t1: stage start time
			// 	t2: stage end time
			for i, t := range ts {
				i, t := i, t

				var t2 time.Time
				if i == L-1 {
					t2 = t
				} else {
					t2 = ts[i+1]
				}

				sl.Add(t, func() {
					// fire
					evtmgr.Fire(
						fmt.Sprintf("cal->%s", v.Type),                        // cal->'type'
						ft.abs, v.ModName, v.StageList[i].Stage, ft.t0, t, t2, // abs, name, stage, t0, t1, t2
					)

					// sched next if this is the last stage
					if i == L-1 {
						if ft.inc() {
							sched_next(ft)
						}
					}
				})
			}

			// start stage line
			sl.Start()
		}

		// start sched
		sched_next(ft)
	}
}

// ============================================================================

func parse_timeline_time(v string) (lstart time.Time, ldays int32, lstop time.Time) {
	/*
		o!@2 3:4:5/6/@100	->	svr-open  time-line ! loop-start / loop-interval / loop-stop
		m!@2 3:4:5/6/@100	->	svr-merge time-line ! loop-start / loop-interval / loop-stop
		a!Ts/6/Ts			->  absolute  time-line ! loop-start / loop-interval / loop-stop
	*/

	v = strings.Trim(v, " ")
	if len(v) < 2 {
		core.Panic("invalid time-line time:", v)
	}

	// get ref time
	var ref time.Time

	merged := worlddata.GetSvrMergeCnt() > 0
	prefix := v[:2]

	if prefix == "o!" {
		if merged {
			return
		} else {
			ref = worlddata.GetSvrCreateTs()
		}
	} else if prefix == "m!" {
		if merged {
			ref = worlddata.GetSvrMergeTs()
		} else {
			return
		}
	} else if prefix == "a!" {
		// do nothing
	} else {
		core.Panic("invalid time-line time:", v)
	}

	// parse
	v = v[2:]
	arr := strings.Split(v, "/")
	L := len(arr)

	// loop start ts
	if ref.IsZero() { // abs-time
		lstart = core.ParseTime(arr[0])
	} else {
		lstart = core.ParseRelativeTime(ref, arr[0])
	}

	// loop days
	if L > 1 {
		ldays = core.Atoi32(arr[1])
		if ldays <= 0 {
			core.Panic("invalid time-line time:", v)
		}
	}

	// loop stop ts
	if L > 2 {
		if ref.IsZero() { // abs-time
			lstop = core.ParseTime(arr[2])
		} else {
			lstop = core.ParseRelativeTime(ref, arr[2])
		}
	}

	// adjust lstart if looped
	if ldays > 0 {
		now := time.Now()

		// if history, advanced to nearby
		if lstart.Before(now) {
			days := int32(now.Sub(lstart).Hours() / 24)
			lstart = lstart.AddDate(0, 0, int(days/ldays*ldays))
		}
	}

	return
}

// ============================================================================

func parse_first_time(v string) *first_time_t {
	v = strings.Trim(v, " ")
	if len(v) < 2 {
		core.Panic("invalid calendar first time:", v)
	}

	var ret first_time_t

	if v[1] == '!' {
		ret.t0, ret.ldays, ret.lstop = parse_timeline_time(v)
		if ret.t0.IsZero() || !ret.lstop.IsZero() && !ret.t0.Before(ret.lstop) {
			return nil
		}
		ret.abs = v[0] != 'o' && v[0] != 'm'
	} else {
		ret.unit, ret.t0 = core.ParseRepeatTime(v)
		ret.abs = true
	}

	return &ret
}

func (self *first_time_t) inc() bool {
	var t time.Time

	if self.unit == "" {
		if self.ldays <= 0 {
			return false
		}

		t = self.t0.AddDate(0, 0, int(self.ldays))
		if !self.lstop.IsZero() && !t.Before(self.lstop) {
			return false
		}
	} else {
		t = core.AddTimeByUnit(self.t0, self.unit, 1)
	}

	// ok
	self.t0 = t
	return true
}
