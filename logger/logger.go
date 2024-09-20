package logger

type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Trace(arg interface{})
	Debug(arg interface{})
	Info(arg interface{})
	Warn(arg interface{})
	Error(arg interface{})
	Fatal(arg interface{})
	Panic(arg interface{})
}
