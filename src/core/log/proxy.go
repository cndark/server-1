package log

import (
	"github.com/op/go-logging"
)

// ============================================================================

var (
	logger = logging.MustGetLogger("")
)

// ============================================================================

func init() {
	logger.ExtraCalldepth = 1
}

// ============================================================================

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Notice(args ...interface{}) {
	logger.Notice(args...)
}

func Noticef(format string, args ...interface{}) {
	logger.Noticef(format, args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Critical(args ...interface{}) {
	logger.Critical(args...)
}

func Criticalf(format string, args ...interface{}) {
	logger.Criticalf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
