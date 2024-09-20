package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type logger struct {
	*logrus.Entry
}

func NewLogger(instance string) Logger {
	log := &logger{
		logrus.WithFields(logrus.Fields{
			// ep: entrypoint
			"ep": strings.ToLower(instance),
		}),
	}

	log.Logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})

	log.Logger.SetLevel(logrus.DebugLevel)

	return log
}

func (l *logger) logf(level logrus.Level, format string, args ...interface{}) {
	l.Entry.Logf(level, format, args...)
}

func (l *logger) log(level logrus.Level, arg interface{}) {
	l.Entry.Log(level, arg)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.logf(logrus.TraceLevel, format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logf(logrus.InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logf(logrus.WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logf(logrus.FatalLevel, format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.logf(logrus.PanicLevel, format, args...)
}

func (l *logger) Trace(arg interface{}) {
	l.log(logrus.TraceLevel, arg)
}

func (l *logger) Debug(arg interface{}) {
	l.log(logrus.DebugLevel, arg)
}

func (l *logger) Info(arg interface{}) {
	l.log(logrus.InfoLevel, arg)
}

func (l *logger) Warn(arg interface{}) {
	l.log(logrus.WarnLevel, arg)
}

func (l *logger) Error(arg interface{}) {
	l.log(logrus.ErrorLevel, arg)
}

func (l *logger) Fatal(arg interface{}) {
	l.log(logrus.FatalLevel, arg)
}

func (l *logger) Panic(arg interface{}) {
	l.log(logrus.PanicLevel, arg)
}
