package core

import (
	"regexp"
	"strings"
	"time"
)

// ============================================================================

var (
	re_datetime     = regexp.MustCompile(`^((\d+)\-(\d+)\-(\d+))?\s*((\d+)(:(\d+)(:(\d+))?)?)?$`)
	re_datetime_sep = regexp.MustCompile(`[-\s:]`)
)

// ============================================================================

func ParseTime(v string) (t time.Time) {
	// format:
	//	[Y-m-d] [H:[M:[S]]]

	v = strings.Trim(v, " ")

	p := re_datetime.FindStringSubmatch(v)
	if p == nil {
		Panic("invalid date:", v)
	}

	y, m, d := int(Atoi32(p[2])), time.Month(Atoi32(p[3])), int(Atoi32(p[4]))
	H, M, S := int(Atoi32(p[6])), int(Atoi32(p[8])), int(Atoi32(p[10]))

	if y == 0 && m == 0 && d == 0 {
		y, m, d = time.Now().Date()
	}

	return time.Date(y, m, d, H, M, S, 0, time.Local)
}

func ParseRelativeTime(ref time.Time, v string) (t time.Time) {
	// format:
	//	* @2 10:50:0	@相对ref当天0点的  天数 时:分:秒
	//	* +2 10:50:0	+相对ref的        天数 时:分:秒

	v = strings.Trim(v, " ")

	if !strings.HasPrefix(v, "@") && !strings.HasPrefix(v, "+") {
		Panic("invalid relative-time:", v)
	}

	// prefix and content
	prefix := v[:1]
	e := v[1:]

	// relative to when
	if prefix == "@" {
		ref = StartOfDay(ref)
	}

	// relative amount
	var p []int
	for _, v := range re_datetime_sep.Split(e, -1) {
		p = append(p, int(Atoi32(v)))
	}
	p = append(p, 0, 0, 0)

	// return
	return ref.
		AddDate(0, 0, p[0]).
		Add(time.Hour * time.Duration(p[1])).
		Add(time.Minute * time.Duration(p[2])).
		Add(time.Second * time.Duration(p[3]))
}

func ParseRepeatTime(v string) (unit string, t time.Time) {
	// format:
	//	unit/time-string

	v = strings.Trim(v, " ")

	ut := strings.Split(v, "/")
	if len(ut) != 2 {
		Panic("invalid repeat-time:", v)
	}

	// unit
	unit = ut[0]

	// time-string
	var p []int
	for _, v := range re_datetime_sep.Split(ut[1], -1) {
		p = append(p, int(Atoi32(v)))
	}
	p = append(p, 0, 0, 0, 0)

	// now
	now := time.Now()
	y, m, d := now.Date()
	H, M, _ := now.Clock()

	switch unit {
	case "M":
		t = time.Date(y, m, d, H, M, p[0], 0, time.Local)

	case "H":
		t = time.Date(y, m, d, H, p[0], p[1], 0, time.Local)

	case "d":
		t = time.Date(y, m, d, p[0], p[1], p[2], 0, time.Local)

	case "m":
		t = time.Date(y, m, p[0], p[1], p[2], p[3], 0, time.Local)

	case "y":
		t = time.Date(y, time.Month(p[0]), p[1], p[2], p[3], p[4], 0, time.Local)

	case "w":
		t = time.Date(y, m, d, p[1], p[2], p[3], 0, time.Local)
		t = t.AddDate(0, 0, p[0]%7-int(t.Weekday()))

	default:
		Panic("invalid repeat-time unit:", unit)
	}

	// return the key-time just has passed
	if !t.Before(now) {
		t = AddTimeByUnit(t, unit, -1)
	}

	return
}

func AddTimeByUnit(t time.Time, unit string, n int) time.Time {
	switch unit {
	case "S":
		return t.Add(time.Second * time.Duration(n))

	case "M":
		return t.Add(time.Minute * time.Duration(n))

	case "H":
		return t.Add(time.Hour * time.Duration(n))

	case "d":
		return t.AddDate(0, 0, n)

	case "m":
		return t.AddDate(0, n, 0)

	case "y":
		return t.AddDate(n, 0, 0)

	case "w":
		return t.AddDate(0, 0, 7*n)

	default:
		Panic("invalid time unit:", unit)
		return t
	}
}

func StartOfDay(t time.Time) time.Time {
	y, M, d := t.Date()
	return time.Date(y, M, d, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	y, M, d := t.Date()
	return time.Date(y, M, d, 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
}

func IsSameDay(t1, t2 time.Time) bool {
	y1, M1, d1 := t1.Date()
	y2, M2, d2 := t2.Date()

	return y1 == y2 && M1 == M2 && d1 == d2
}

func IsSameWeek(t1, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()

	return y1 == y2 && w1 == w2
}

func IsSameMonth(t1, t2 time.Time) bool {
	y1, M1, _ := t1.Date()
	y2, M2, _ := t2.Date()

	return y1 == y2 && M1 == M2
}

func MinTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	} else {
		return b
	}
}

// 今天天数相对key的天数,当天为1
func DistanceDays(key time.Time) int32 {
	key_0 := StartOfDay(key)
	now_0 := StartOfDay(time.Now())

	day := int32(now_0.Sub(key_0).Hours()) / 24
	return day + 1
}
