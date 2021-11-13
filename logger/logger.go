package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type LogLevel = logrus.Level

type logger struct {
	ServiceName string
	Logger      *logrus.Logger
	Entry       *logrus.Entry
}

type Logger interface {
	Info(opts ...interface{})
	Warn(opts ...interface{})
	Debug(opts ...interface{})
	Error(opts ...interface{})
	NewChild(cname string) Logger
	SetLevel(level string)
}

func getFields(opts ...interface{}) (logrus.Fields, []interface{}) {
	fields := logrus.Fields{}
	var interfaces []interface{}

	if opts != nil {
		for _, values := range opts {
			switch val := values.(type) {
			case map[string]interface{}:
				for k, v := range val {
					fields[k] = v
				}
			case interface{}:
				interfaces = append(interfaces, values)
			}
		}
	}

	return fields, interfaces
}

func (l logger) Warn(opts ...interface{}) {
	fields, ms := getFields(opts)
	l.Entry.WithFields(fields).Warn(ms...)
}

func (l logger) Info(opts ...interface{}) {
	fields, ms := getFields(opts)
	l.Entry.WithFields(fields).Info(ms...)
}

func (l logger) Debug(opts ...interface{}) {
	fields, ms := getFields(opts)
	l.Entry.WithFields(fields).Debug(ms...)
}

func (l logger) Error(opts ...interface{}) {
	fields, ms := getFields(opts)
	l.Entry.WithFields(fields).Error(ms...)
}

func (l logger) NewChild(cname string) Logger {
	l.Entry = l.Entry.WithFields(logrus.Fields{
		"childService": cname,
	})

	return l
}

func (l logger) SetLevel(level string) {
	switch level {
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
}

func New(serviceName string) Logger {
	l := logger{
		ServiceName: serviceName,
	}

	l.Logger = logrus.New()
	l.Entry = l.Logger.WithFields(logrus.Fields{
		"service": l.ServiceName,
	})

	format := os.Getenv("LOG_FORMAT")
	if format == "console" {
		l.Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		l.Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	logLevel := os.Getenv("LOG_LEVEL")
	l.SetLevel(logLevel)

	return l
}
