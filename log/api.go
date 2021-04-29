package log

type Level string

const (
	DebugLevel  Level = "debug"
	InfoLevel   Level = "info"
	WarnLevel   Level = "waring"
	ErrorLevel  Level = "error"
	DPanicLevel Level = "dPanic"
	PanicLevel  Level = "panic"
	FatalLevel  Level = "fatal"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Panic(v ...interface{})
	Fatal(v ...interface{})
}

var globalLogger Logger

func Debug(v ...interface{}) {
	globalLogger.Debug(v...)
}
func Info(v ...interface{}) {
	globalLogger.Info(v...)
}
func Warn(v ...interface{}) {
	globalLogger.Warn(v...)
}
func Error(v ...interface{}) {
	globalLogger.Error(v...)
}
func Panic(v ...interface{}) {
	globalLogger.Panic(v...)
}
func Fatal(v ...interface{}) {
	globalLogger.Fatal(v...)
}

func SetLogger(l Logger) {
	globalLogger = l
}

func Global() Logger {
	return globalLogger
}
