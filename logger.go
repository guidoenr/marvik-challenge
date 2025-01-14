// logger.go
package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// InitLogger initializes the logger.
func InitLogger() zerolog.Logger {
	logLevel := zerolog.InfoLevel

	if os.Getenv("DEBUG") == "true" {
		logLevel = zerolog.DebugLevel
	}

	// human-readable output (plain text)
	consoleWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
		w.NoColor = false
	})

	return zerolog.New(consoleWriter).With().Timestamp().Logger().Level(logLevel)
}
