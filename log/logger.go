package log

import (
	"log"
	"os"
)

type IntLogLevel int8

const (
	IntLogDebugLevel IntLogLevel = iota - 1
	IntLogInfoLevel
	IntLogWarnLevel
	IntLogErrorLevel
	IntLogPanicLevel
	IntLogFatalLevel
)

func UnmarshalLevel(l Level) IntLogLevel {
	switch l {
	case DebugLevel:
		return IntLogDebugLevel
	case InfoLevel:
		return IntLogInfoLevel
	case WarnLevel:
		return IntLogWarnLevel
	case ErrorLevel:
		return IntLogErrorLevel
	case PanicLevel:
		return IntLogPanicLevel
	case FatalLevel:
		return IntLogFatalLevel
	default:
		return IntLogInfoLevel
	}
}

type logger struct {
	*log.Logger
	Service     string
	intLogLevel IntLogLevel
}

func NewLogger(cfg *LogConfig) logger {
	file, err := os.Create(cfg.Filename)
	if err != nil {
		log.Fatalln("fail to create log file!", err)
	}
	goLogger := log.New(file, "", log.LstdFlags|log.Llongfile)

	goLogger.SetFlags(log.LstdFlags)

	lg := logger{
		Logger:  goLogger,
		Service: cfg.Service,
	}
	lg.SetLevel(cfg.Level)
	return lg
}

func (l *logger) SetLevel(lv Level) {
	l.intLogLevel = UnmarshalLevel(lv)
}

func (l *logger) Enabled(lv IntLogLevel) bool {
	return lv >= l.intLogLevel
}

func (l *logger) Debug(v ...interface{}) {
	if !l.Enabled(IntLogDebugLevel) {
		return
	}
	l.Logger.Println(v...)
}

func (l *logger) Info(v ...interface{}) {
	if !l.Enabled(IntLogInfoLevel) {
		return
	}
	l.Logger.Println(v...)
}

func (l *logger) Warn(v ...interface{}) {
	if !l.Enabled(IntLogWarnLevel) {
		return
	}
	l.Logger.Println(v...)
}

func (l *logger) Error(v ...interface{}) {
	if !l.Enabled(IntLogErrorLevel) {
		return
	}
	l.Logger.Println(v...)
}

func (l *logger) Panic(v ...interface{}) {
	if !l.Enabled(IntLogPanicLevel) {
		return
	}
	l.Logger.Panic(v...)
}

func (l *logger) Fatal(v ...interface{}) {
	if !l.Enabled(IntLogFatalLevel) {
		return
	}
	l.Logger.Fatal(v...)
}
