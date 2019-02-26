package log

import (
	"fmt"

	"github.com/pkg/errors"
)

func NewTracingLogger(next Logger) Logger {
	return &tracingLogger{
		next: next,
	}
}

type stackTracer interface {
	error
	StackTrace() errors.StackTrace
}

type tracingLogger struct {
	next Logger
}

func (l *tracingLogger) Trace(msg string, keyvals ...interface{}) {
	l.next.Trace(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Debug(msg string, keyvals ...interface{}) {
	l.next.Debug(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Info(msg string, keyvals ...interface{}) {
	l.next.Info(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Warn(msg string, keyvals ...interface{}) {
	l.next.Warn(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Error(msg string, keyvals ...interface{}) {
	l.next.Error(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Fatal(msg string, keyvals ...interface{}) {
	l.next.Fatal(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) With(keyvals ...interface{}) Logger {
	return &tracingLogger{next: l.next.With(formatErrors(keyvals)...)}
}

func formatErrors(keyvals []interface{}) []interface{} {
	newKeyvals := make([]interface{}, len(keyvals))
	copy(newKeyvals, keyvals)
	for i := 0; i < len(newKeyvals)-1; i += 2 {
		if err, ok := newKeyvals[i+1].(stackTracer); ok {
			newKeyvals[i+1] = tracedError{err}
		}
	}
	return newKeyvals
}

type tracedError struct {
	wrapped stackTracer
}

var _ stackTracer = tracedError{}

func (t tracedError) StackTrace() errors.StackTrace {
	return t.wrapped.StackTrace()
}

func (t tracedError) Cause() error {
	return t.wrapped
}

func (t tracedError) Error() string {
	return fmt.Sprintf("%+v", t.wrapped)
}