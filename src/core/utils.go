package core

import (
	"fmt"
)

// ============================================================================

func DefTrue(v []bool) bool {
	return len(v) == 0 || v[0]
}

func DefFalse(v []bool) bool {
	return len(v) > 0 && v[0]
}

func Panic(v ...interface{}) {
	panic(fmt.Sprintln(v...))
}
