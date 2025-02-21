package log

import (
	"io"

	"github.com/rs/zerolog"
)

// ModuleKey defines a module logging key.
const ModuleKey = "module"

// ContextKey is used to store the logger in the context.
var ContextKey struct{}

// Logger is the PellApp SDK logger interface.
// It maintains as much backward compatibility with the PellDVS logger as possible.
// All functionalities of the logger are available through the Impl() method.
type Logger interface {
	// Info takes a message and a set of key/value pairs and logs with level INFO.
	// The key of the tuple must be a string.
	Info(msg string, keyVals ...any)

	// Error takes a message and a set of key/value pairs and logs with level ERR.
	// The key of the tuple must be a string.
	Error(msg string, keyVals ...any)

	// Debug takes a message and a set of key/value pairs and logs with level DEBUG.
	// The key of the tuple must be a string.
	Debug(msg string, keyVals ...any)

	// With returns a new wrapped logger with additional context provided by a set.
	With(keyVals ...any) Logger

	// Impl returns the underlying logger implementation.
	// It is used to access the full functionalities of the underlying logger.
	// Advanced users can type cast the returned value to the actual logger.
	Impl() any
}

type zeroLogWrapper struct {
	*zerolog.Logger
}

// NewLogger returns a new logger that writes to the given destination.
//
// Typical usage from a main function is:
//
//	logger := log.NewLogger(os.Stderr)
//
// Stderr is the typical destination for logs,
// so that any output from your application can still be piped to other processes.
func NewLogger(dst io.Writer, options ...Option) Logger {
	logCfg := defaultConfig
	for _, opt := range options {
		opt(&logCfg)
	}

	output := dst
	if !logCfg.OutputJSON {
		output = zerolog.ConsoleWriter{
			Out:        dst,
			NoColor:    !logCfg.Color,
			TimeFormat: logCfg.TimeFormat,
		}
	}

	if logCfg.Filter != nil {
		output = NewFilterWriter(output, logCfg.Filter)
	}

	logger := zerolog.New(output)
	if logCfg.TimeFormat != "" {
		logger = logger.With().Timestamp().Logger()
	}

	if logCfg.Level != zerolog.NoLevel {
		logger = logger.Level(logCfg.Level)
	}

	return zeroLogWrapper{&logger}
}

// NewCustomLogger returns a new logger with the given zerolog logger.
func NewCustomLogger(logger zerolog.Logger) Logger {
	return zeroLogWrapper{&logger}
}

// Info takes a message and a set of key/value pairs and logs with level INFO.
// The key of the tuple must be a string.
func (l zeroLogWrapper) Info(msg string, keyVals ...interface{}) {
	l.Logger.Info().Fields(keyVals).Msg(msg)
}

// Error takes a message and a set of key/value pairs and logs with level DEBUG.
// The key of the tuple must be a string.
func (l zeroLogWrapper) Error(msg string, keyVals ...interface{}) {
	l.Logger.Error().Fields(keyVals).Msg(msg)
}

// Debug takes a message and a set of key/value pairs and logs with level ERR.
// The key of the tuple must be a string.
func (l zeroLogWrapper) Debug(msg string, keyVals ...interface{}) {
	l.Logger.Debug().Fields(keyVals).Msg(msg)
}

// With returns a new wrapped logger with additional context provided by a set.
func (l zeroLogWrapper) With(keyVals ...interface{}) Logger {
	logger := l.Logger.With().Fields(keyVals).Logger()
	return zeroLogWrapper{&logger}
}

// Impl returns the underlying zerolog logger.
// It can be used to used zerolog structured API directly instead of the wrapper.
func (l zeroLogWrapper) Impl() interface{} {
	return l.Logger
}
