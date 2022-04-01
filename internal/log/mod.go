// Package log implements wrapper around logging libraries.
package log

import "log"

// Info logs an informational log message.
func Info(message string) {
	oldPrefix := log.Prefix()
	log.SetPrefix("INFO ")
	log.Printf("%s", message)
	log.SetPrefix(oldPrefix)
}

// Error logs an error message and the exception thrown.
func Error(message string, err error) {
	oldPrefix := log.Prefix()
	log.SetPrefix("ERROR ")
	log.Printf("%s: %v", message, err)
	log.SetPrefix(oldPrefix)
}

// Fatal causes the program to stop with error message.
func Fatal(err interface{}) {
	log.Fatal(err)
}
