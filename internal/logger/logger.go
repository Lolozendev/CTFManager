package logger

import (
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	instance *log.Logger
	logFile  *os.File
	once     sync.Once
)

// initLogger initializes the logger singleton
func initLogger() {
	var err error

	// Open single log file
	logFile, err = os.OpenFile(
		"ctfmanager.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// Create logger with charmbracelet's pretty terminal output
	instance = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05",
		ReportCaller:    false,
	})

	// Write to both terminal (with colors) and file (plain text)
	instance.SetOutput(io.MultiWriter(os.Stderr, logFile))
}

// Get returns the global logger instance
func Get() *log.Logger {
	once.Do(initLogger)
	return instance
}

// Close closes the log file - call this on graceful shutdown
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
