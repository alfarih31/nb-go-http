package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
)

type LogLevel = logrus.Level

type logger struct {
	ServiceName string
	Logger      *logrus.Logger
	Entry       *logrus.Entry
	Level       LogLevel
}

type Logger interface {
	Info(m interface{}, opts ...interface{}) Logger
	Warn(m interface{}, opts ...interface{}) Logger
	Debug(m interface{}, opts ...interface{}) Logger
	Error(m interface{}, opts ...interface{}) Logger
	NewChild(cname string) Logger
	SetLevel(level string) Logger
	AddHook(hook logrus.Hook) Logger
	GetLevel() LogLevel
}

func getFields(m interface{}, opts []interface{}) (logrus.Fields, []interface{}) {
	fields := logrus.Fields{}
	interfaces := []interface{}{m}

	if opts != nil {
		for _, values := range opts {
			if values == nil {
				continue
			}

			tValues := reflect.ValueOf(values)
			switch tValues.Kind() {
			case reflect.Map:
				for _, k := range tValues.MapKeys() {
					fields[k.String()] = tValues.MapIndex(k).Interface()
				}
			default:
				interfaces = append(interfaces, values)
			}
		}
	}

	return fields, interfaces
}

func (l *logger) Warn(m interface{}, opts ...interface{}) Logger {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Warn(ms...)

	return l
}

func (l *logger) Info(m interface{}, opts ...interface{}) Logger {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Info(ms...)

	return l
}

func (l *logger) Debug(m interface{}, opts ...interface{}) Logger {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Debug(ms...)

	return l
}

func (l *logger) Error(m interface{}, opts ...interface{}) Logger {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Error(ms...)

	return l
}

func (l *logger) NewChild(cname string) Logger {
	l.Entry = l.Entry.WithFields(logrus.Fields{
		"childService": cname,
	})

	return l
}

func (l *logger) SetLevel(level string) Logger {
	switch level {
	case "debug":
		l.Logger.SetLevel(logrus.DebugLevel)
		l.Level = logrus.DebugLevel

	case "error":
		l.Logger.SetLevel(logrus.ErrorLevel)
		l.Level = logrus.ErrorLevel

	case "info":
		l.Logger.SetLevel(logrus.InfoLevel)
		l.Level = logrus.InfoLevel

	case "warn":
		l.Logger.SetLevel(logrus.WarnLevel)
		l.Level = logrus.WarnLevel

	default:
		l.Logger.SetLevel(logrus.InfoLevel)
		l.Level = logrus.InfoLevel

	}

	return l
}

func (l *logger) AddHook(hook logrus.Hook) Logger {
	l.Logger.AddHook(hook)

	return l
}

func (l *logger) GetLevel() LogLevel {
	return l.Level
}

func New(serviceName string) Logger {
	l := new(logger)
	l.ServiceName = serviceName

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
