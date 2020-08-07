package log

import (
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

	outType = strings.ToLower(outType)

	switch outType {
	case "json":
		logger.SetFormatter(new(logrus.JSONFormatter))
	case "text":
		logger.SetFormatter(new(logrus.TextFormatter))
	default:
		logger.SetFormatter(new(logrus.TextFormatter))
	}
}

// Fatal wraps log.Fatal
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf wraps log.Fatalf
func Fatalf(fmt string, args ...interface{}) {
	logger.Fatalf(fmt, args...)
}

// Debug wraps log.Debug
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf wraps log.Debugf
func Debugf(fmt string, args ...interface{}) {
	logger.Debugf(fmt, args...)
}

// Info wraps log.Info
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof wraps log.Infof
func Infof(fmt string, args ...interface{}) {
	logger.Infof(fmt, args...)
}

// Infof wraps log.Infof
func InfoWithFields(fields logrus.Fields, msg string) {
	logger.WithFields(fields).Info(msg)
}

// Warn wraps log.Warn
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf wraps log.Warnf
func Warnf(fmt string, args ...interface{}) {
	logger.Warnf(fmt, args...)
}

// Error wraps log.Error
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf wraps log.Errorf
func Errorf(fmt string, args ...interface{}) {
	logger.Errorf(fmt, args...)
}
