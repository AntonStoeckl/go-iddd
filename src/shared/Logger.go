package shared

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func NewStandardLogger() *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := &Logger{zerolog.New(output).With().Caller().Timestamp().Logger()}

	return logger
}

func NewNilLogger() *Logger {
	logger := &Logger{zerolog.New(nil)}

	return logger
}

// Verbose() is needed to satisfy the Logger interface of github.com/golang-migrate/migrate/v4
func (logger *Logger) Verbose() bool {
	return true
}
