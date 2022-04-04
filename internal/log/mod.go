// Package log implements wrapper around logging libraries.
package log

import "log"

// Info logs an informational log message.
func Info(message string) {
	logPrefixf("INFO", message)
}

// Warn logs a warning log message.
func Warn(message string) {
	logPrefixf("WARN", message)
}

// Error logs an error message and the exception thrown.
func Error(message string, err error) {
	logPrefixf("ERROR", "%s: %v", message, err)
}

func logPrefixf(prefix string, message string, v ...any) {
	oldPrefix := log.Prefix()
	log.SetPrefix(prefix + " ")
	log.Printf(message, v...)
	log.SetPrefix(oldPrefix)
}

// Fatal causes the program to stop with error message.
func Fatal(err interface{}) {
	log.Fatal(err)
}
