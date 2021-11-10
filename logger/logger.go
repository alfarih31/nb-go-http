package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type TLogger struct {
	ServiceName string
	Logger      *logrus.Logger
	Entry       *logrus.Entry
}

type ILogger interface {
	Info(m interface{}, opts map[string]interface{})
	Warn(m interface{}, opts map[string]interface{})
	Debug(m interface{}, opts map[string]interface{})
	Error(m interface{}, opts map[string]interface{})
	NewChild(cname string) ILogger
}

func getFields(opts map[string]interface{}) logrus.Fields {
	fields := logrus.Fields{}

	if opts != nil {
		for key, val := range opts {
			fields[key] = val
		}
	}

	return fields
}

func (l TLogger) Warn(m interface{}, opts map[string]interface{}) {
	l.Entry.WithFields(getFields(opts)).Warn(m)
}

func (l TLogger) Info(m interface{}, opts map[string]interface{}) {
	l.Entry.WithFields(getFields(opts)).Info(m)
}

func (l TLogger) Debug(m interface{}, opts map[string]interface{}) {
	l.Entry.WithFields(getFields(opts)).Debug(m)
}

func (l TLogger) Error(m interface{}, opts map[string]interface{}) {
	l.Entry.WithFields(getFields(opts)).Error(m)
}

func (l TLogger) NewChild(cname string) ILogger {
	l.Entry = l.Entry.WithFields(logrus.Fields{
		"childService": cname,
	})

	return l
}

func Logger(serviceName string) ILogger {
	l := TLogger{
		ServiceName: serviceName,
	}

	l.Logger = logrus.New()
	l.Entry = l.Logger.WithFields(logrus.Fields{
		"service": l.ServiceName,
	})

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		l.Logger.SetLevel(logrus.DebugLevel)

	case "error":
		l.Logger.SetLevel(logrus.ErrorLevel)

	case "info":
		l.Logger.SetLevel(logrus.InfoLevel)

	case "warn":
		l.Logger.SetLevel(logrus.WarnLevel)

	default:
		l.Logger.SetLevel(logrus.InfoLevel)

	}

	format := os.Getenv("LOG_FORMAT")
	if format == "console" {
		l.Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		l.Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return l
}
