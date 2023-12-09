package logger

import (
	"github.com/sirupsen/logrus"
)

// declare global logger variable to be used in the whole project
var CLogger = NewLogger()

// Logger defines the logger interface that is exposed via RequestScope.
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Success(args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Successf(format string, args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Debugln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Traceln(args ...interface{})
	Successln(args ...interface{})
}

type logger struct {
	*logrus.Logger
}

// InitLogger initializes the logger instance.
func InitLogger() Logger {
	logger := NewLogger()
	return logger
}

// NewLogger creates a new logger instance.
func NewLogger() Logger {
	return &logger{
		Logger: logrus.New(),
	}
}

// Info logs a message at level Info on the standard logger.
func (l *logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *logger) Trace(args ...interface{}) {
	l.Logger.Trace(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *logger) Success(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.Logger.Tracef(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *logger) Successf(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.Logger.Infoln(args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.Logger.Warnln(args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.Logger.Debugln(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.Logger.Errorln(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.Logger.Fatalln(args...)
}

func (l *logger) Traceln(args ...interface{}) {
	l.Logger.Traceln(args...)
}

func (l *logger) Successln(args ...interface{}) {
	l.Logger.Infoln(args...)
}
