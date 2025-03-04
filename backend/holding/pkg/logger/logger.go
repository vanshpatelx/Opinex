package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/mattn/go-colorable" // Ensures colors work on Windows too
)

// InitLogger initializes zerolog with both console (colored) and file (JSON) logging
func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	// Console output (colored logs)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        colorable.NewColorableStdout(),
		TimeFormat: "15:04:05", // HH:mm:ss format
		NoColor:    false,      // Ensures color output
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("💬 %v", i) // Adds emoji to message
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%v=", i) // Keeps field names visible
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%v", i) // Keeps values visible
		},
	}

	// File output (JSON logs)
	file, err := os.OpenFile("logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	// Multi-writer: console (colored) + file (JSON)
	multiWriter := zerolog.MultiLevelWriter(consoleWriter, file)

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
}

// Info logs an informational message
func Info(msg string, fields ...interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...interface{}) {
	log.Warn().Fields(fields).Msg(msg)
}

// Error logs an error message
func Error(msg string, fields ...interface{}) {
	log.Error().Fields(fields).Msg(msg)
}

// Debug logs a debug message
func Debug(msg string, fields ...interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

// Fatal logs a critical error and exits
func Fatal(msg string, fields ...interface{}) {
	log.Fatal().Fields(fields).Msg(msg)
	os.Exit(1) // Ensure the application exits on fatal errors
}
