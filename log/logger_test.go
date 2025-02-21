package log

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(os.Stdout, ColorOption(true))
	logger.Debug("hello world")

	levelInfoLogger := NewLogger(os.Stdout, LevelOption(zerolog.InfoLevel))
	levelInfoLogger.Debug("this log line should not be displayed, because the log level is set to info")
	levelInfoLogger.Info("this log line should be displayed")
}
