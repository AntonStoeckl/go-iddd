package cmd

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewStandardLogger() *Logger {
	logger := &Logger{logrus.New()}

	logger.SetFormatter(
		&logrus.TextFormatter{
			FullTimestamp: true,
		},
	)

	return logger
}

// Verbose() is needed to satisfy the Logger interface of github.com/golang-migrate/migrate/v4
func (logger *Logger) Verbose() bool {
	return true
}
