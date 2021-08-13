package core

import (
	"strconv"
)

func A16toi32(v string) int32 {
	n, err := strconv.ParseInt(v, 16, 32)
	if err == nil {
		return int32(n)
	} else {
		return 0
	}
}

func I32toa(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func Atoi32(v string) int32 {
	n, err := strconv.ParseInt(v, 10, 32)
	if err == nil {
		return int32(n)
	} else {
		return 0
	}
}

func I64toa(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Atoi64(v string) int64 {
	n, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return n
	} else {
		return 0
	}
}

func Atof32(v string) float32 {
	n, err := strconv.ParseFloat(v, 32)
	if err == nil {
		return float32(n)
	} else {
		return 0
	}
}

func Atof64(v string) float64 {
	n, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return n
	} else {
		return 0
	}
}

func F64toa(v float64) string {
	return strconv.FormatFloat(v, 'f', 4, 64)
}
