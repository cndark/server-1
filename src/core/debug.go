package core

import (
	"runtime/debug"
)

func Callstack() string {
	return string(debug.Stack())
}
