package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

// SetupLogger prepares the logger instance
func SetupLogger(level string, outType string) {
	logger.SetLevel(logrus.DebugLevel)
	if logLevel, err := logrus.ParseLevel(level); err == nil {
		logger.SetLevel(logLevel)
	}
	// logger.WithField("a", "b")

	outType = strings.ToLower(outType)

	switch outType {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	case "text":
		logrus.SetFormatter(new(logrus.TextFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// Fatal wraps log.Fatal
func Fatal(args ...interface{}) {
	logger.Fatal(addHostToArgs(args)...)
}

// Fatalf wraps log.Fatalf
func Fatalf(fmt string, args ...interface{}) {
	logger.Fatalf(addHostToFmt(fmt), addHostToArgs(args)...)
}

// Debug wraps log.Debug
func Debug(args ...interface{}) {
	logger.Debug(addHostToArgs(args)...)
}

// Debugf wraps log.Debugf
func Debugf(fmt string, args ...interface{}) {
	logger.Debugf(addHostToFmt(fmt), addHostToArgs(args)...)
}

// Info wraps log.Info
func Info(args ...interface{}) {
	logger.Info(addHostToArgs(args)...)
}

// Infof wraps log.Infof
func Infof(fmt string, args ...interface{}) {
	logger.Infof(addHostToFmt(fmt), addHostToArgs(args)...)
}

// Warn wraps log.Warn
func Warn(args ...interface{}) {
	logger.Warn(addHostToArgs(args)...)
}

// Warnf wraps log.Warnf
func Warnf(fmt string, args ...interface{}) {
	logger.Warnf(addHostToFmt(fmt), addHostToArgs(args)...)
}

// Error wraps log.Error
func Error(args ...interface{}) {
	logger.Error(addHostToArgs(args)...)
}

// Errorf wraps log.Errorf
func Errorf(fmt string, args ...interface{}) {
	logger.Errorf(addHostToFmt(fmt), addHostToArgs(args)...)
}

func addHostToFmt(fmt string) string {
	return "[%s] " + fmt
}

func addHostToArgs(args ...interface{}) []interface{} {
	host, err := os.Hostname()
	if err == nil {
		var newArgs []interface{}
		newArgs = append(newArgs, host)
		for _, v := range args {
			newArgs = append(newArgs, v)
		}
		return newArgs
	} else {
		return args
	}

}
