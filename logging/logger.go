package logging

import "go.uber.org/zap"

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defaultLogger = Logger{logger.Sugar()}
}

// A logger which writes given message to stderr.
//
// It's safe to call all methods on a nil receiver.
type Logger struct {
	inner *zap.SugaredLogger
}

// Global logger instance
var defaultLogger Logger

// Returns default logger suitable for usage
func DefaultLogger() *Logger {
	return &defaultLogger
}

// Adds a variadic number of fields to the logging context. The first value
// will become a key and the second one will become a value.
func (l *Logger) With(args ...interface{}) *Logger {
	if l == nil || l.inner == nil {
		return nil
	}

	inner := l.inner.With(args...)
	return &Logger{inner}
}

// Logs formatted debug message.
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Debugf(template, args...)
	})
}

// Logs formatted info message.
func (l *Logger) Infof(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Infof(template, args...)
	})
}

// Logs formatted warning message.
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Warnf(template, args...)
	})
}

// Logs formatted error message.
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Errorf(template, args...)
	})
}

// Executes given closure only if receiver and inner loggers are not nil.
func (l *Logger) safeExec(f func()) {
	if l != nil && l.inner != nil {
		f()
	}
}
