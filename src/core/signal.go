package core

import (
	"os"
	"os/signal"
	"syscall"
)

func OnSignal(f func(os.Signal)) {
	go func() {
		c := make(chan os.Signal, 16)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.Signal(0xa))

		for s := range c {
			f(s)
			if s != syscall.SIGHUP && s != syscall.Signal(0xa) {
				break
			}
		}
	}()
}
