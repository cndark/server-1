package log

import (
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
)

// ============================================================================

var (
	log_fmt = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} [%{level:.1s}] => %{message} | [%{shortfile:s}] %{color:reset}`,
	)
)

// ============================================================================

var (
	quit chan int

	log_fn    string
	log_h     *os.File
	back_end  *logging.LogBackend
	last_date string
)

// ============================================================================

func Open(fn string) {
	if quit != nil {
		return
	}

	quit = make(chan int)

	log_fn = fn

	// prepare backend
	if log_fn == "" {
		// output to stdout
		log_h = os.Stdout
	} else {
		// output to file
		log_h, _ = os.OpenFile(log_fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0644))
		fi, err := log_h.Stat()
		if err == nil {
			last_date = fi.ModTime().Format("2006-01-02")
		} else {
			last_date = time.Now().Format("2006-01-02")
		}
	}

	// set backend
	back_end = logging.NewLogBackend(log_h, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(back_end, log_fmt),
	)

	// set level
	if os.Getenv("LOG_LEVEL_INFO") == "true" {
		logging.SetLevel(logging.INFO, "")
	} else {
		logging.SetLevel(logging.DEBUG, "")
	}

	// log rotation check
	if log_fn != "" {
		// open check
		check_rotation()

		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer func() {
				ticker.Stop()
				quit <- 0
			}()

			for {
				select {
				case <-ticker.C:
					// runtime check
					check_rotation()

				case <-quit:
					// close check
					check_rotation()
					return
				}
			}
		}()
	}
}

func Close() {
	if quit == nil {
		return
	}

	if log_fn != "" {
		quit <- 0
		<-quit
	}
	close(quit)
	quit = nil

	log_h.Close()
}

// ============================================================================

func check_rotation() {

	new_date := time.Now().Format("2006-01-02")

	if new_date != last_date {
		// rename log
		newpath := fmt.Sprintf("log/%s.%s", log_fn, last_date)
		os.MkdirAll("log", os.FileMode(0755))
		os.Rename(log_fn, newpath)

		// create new log
		old_h := log_h
		log_h, _ = os.Create(log_fn)
		back_end.Logger.SetOutput(log_h)
		old_h.Close()

		// update last date
		last_date = new_date
	}
}
