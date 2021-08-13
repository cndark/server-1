package core

import "fw/src/core/log"

func Go(f func()) {
	go PCall(f)
}

func PCall(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("critical exception:", err)
			log.Error(Callstack())
		}
	}()

	f()
}
