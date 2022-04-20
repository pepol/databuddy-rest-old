// Package log implements wrapper around logging libraries.
package log

import (
	stdlog "log"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger configures zerolog logger correctly, including stdlib log override.
func Logger(level zerolog.Level, devel bool) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var l zerolog.Logger

	if devel {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	} else {
		l = log.Logger
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(l)

	zerolog.SetGlobalLevel(level)
	return l
}

// Debug message wrapper.
func Debug(format string, args ...any) {
	log.Debug().Msgf(format, args...)
}

// Info message wrapper.
func Info(format string, args ...any) {
	log.Info().Msgf(format, args...)
}

// Warn message wrapper.
func Warn(format string, args ...any) {
	log.Warn().Msgf(format, args...)
}

// Error message wrapper.
func Error(message string, err error) {
	log.Error().Msgf("%s: %v", message, err)
}

// Fatal error message wrapper.
func Fatal(err error) {
	stdlog.Fatal(err)
}

// BadgerLogger is a wrapper for Badger v3 to accept this logging library.
type BadgerLogger struct{}

// GetBadgerLogger returns BadgerLogger for this library.
func GetBadgerLogger() *BadgerLogger {
	return &BadgerLogger{}
}

// Debugf logs a debug message.
func (l *BadgerLogger) Debugf(format string, args ...any) {
	Debug(format, args...)
}

// Infof logs an informational message.
func (l *BadgerLogger) Infof(format string, args ...any) {
	Info(format, args...)
}

// Warningf logs a warning message.
func (l *BadgerLogger) Warningf(format string, args ...any) {
	Warn(format, args...)
}

// Errorf logs an error message.
func (l *BadgerLogger) Errorf(format string, args ...any) {
	log.Error().Msgf(format, args...)
}
