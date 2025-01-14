// logger.go
package main

import (
	"os"

	"github.com/rs/zerolog"
)

// InitLogger initializes the global log variable
func InitLogger() zerolog.Logger {

	logLevel := zerolog.InfoLevel

	if os.Getenv("DEBUG") == "true" {
		logLevel = zerolog.DebugLevel
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)
}
