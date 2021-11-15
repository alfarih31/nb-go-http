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
}

type Logger interface {
	Info(m interface{}, opts ...interface{})
	Warn(m interface{}, opts ...interface{})
	Debug(m interface{}, opts ...interface{})
	Error(m interface{}, opts ...interface{})
	NewChild(cname string) Logger
	SetLevel(level string)
	AddHook(hook logrus.Hook)
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

func (l logger) Warn(m interface{}, opts ...interface{}) {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Warn(ms...)
}

func (l logger) Info(m interface{}, opts ...interface{}) {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Info(ms...)
}

func (l logger) Debug(m interface{}, opts ...interface{}) {
	fields, ms := getFields(m, opts)
	l.Entry.WithFields(fields).Debug(ms...)
}

func (l logger) Error(m interface{}, opts ...interface{}) {
	fields, ms := getFields(m, opts)
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

func (l logger) AddHook(hook logrus.Hook) {
	l.Logger.AddHook(hook)
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
