package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Set logger output to stdout
	logger.SetOutput(os.Stdout)

	// Set default logging level
	logger.SetLevel(logrus.InfoLevel)

	// Set JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger
}
