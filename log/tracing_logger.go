package log

import (
	"fmt"

	"github.com/pkg/errors"
)

// NewTracingLogger enables tracing by wrapping all errors (if they
// implement stackTracer interface) in tracedError.
//
// All errors returned by https://github.com/pkg/errors implement stackTracer
// interface.
//
// For debugging purposes only as it doubles the amount of allocations.
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

func (l *tracingLogger) Impl() any {
	return l
}

func (l *tracingLogger) Info(msg string, keyvals ...any) {
	l.next.Info(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Debug(msg string, keyvals ...any) {
	l.next.Debug(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Warn(msg string, keyvals ...any) {
	l.next.Warn(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) Error(msg string, keyvals ...any) {
	l.next.Error(msg, formatErrors(keyvals)...)
}

func (l *tracingLogger) With(keyvals ...any) Logger {
	return &tracingLogger{next: l.next.With(formatErrors(keyvals)...)}
}

func formatErrors(keyvals []any) []any {
	newKeyvals := make([]any, len(keyvals))
	copy(newKeyvals, keyvals)
	for i := 0; i < len(newKeyvals)-1; i += 2 {
		if err, ok := newKeyvals[i+1].(stackTracer); ok {
			newKeyvals[i+1] = tracedError{err}
		}
	}
	return newKeyvals
}

// tracedError wraps a stackTracer and just makes the Error() result
// always return a full stack trace.
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
