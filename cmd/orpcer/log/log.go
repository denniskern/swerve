package log

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func SetupLogger(level string, outType string) *logrus.Logger {
	logger := logrus.New()
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
	return logger
}
