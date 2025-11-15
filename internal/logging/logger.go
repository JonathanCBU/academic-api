package logging

import (
	"maps"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetLogLevel() logrus.Level {
	logLevel := os.Getenv(EnvLogLevel)
	switch strings.ToLower(logLevel) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	}
	return logrus.InfoLevel
}

type HttpLog struct {
	Err    error
	Msg    string
	Level  logrus.Level
	Fields logrus.Fields
}

type LoggingContext struct {
	Formatter logrus.Formatter
	Level     logrus.Level
}

func NewLoggingContext() *LoggingContext {
	logFormat := os.Getenv(EnvLogFormat)
	switch logFormat {
	case "json":
		return &LoggingContext{
			Formatter: &logrus.JSONFormatter{},
			Level:     GetLogLevel(),
		}
	case "text":
		return &LoggingContext{
			Formatter: &logrus.TextFormatter{},
			Level:     GetLogLevel(),
		}
	default:
		return &LoggingContext{
			Formatter: &logrus.JSONFormatter{},
			Level:     GetLogLevel(),
		}
	}
}

func (lc *LoggingContext) Log(log *HttpLog) {
	logData := logrus.Fields{
		"error": log.Err,
		"msg":   log.Msg,
	}
	if log.Fields != nil {
		maps.Copy(logData, log.Fields)
	}
	switch log.Level {
	case logrus.FatalLevel:
		logrus.WithFields(logData).Fatal(log.Msg, log.Err)
	case logrus.PanicLevel:
		logrus.WithFields(logData).Panic(log.Msg, log.Err)
	default:
		logrus.WithFields(logData).Log(log.Level, log.Msg, log.Err)
	}
}
