package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	instance   *log.Logger
	eventsFile *os.File
	errorsFile *os.File
	once       sync.Once
)

// levelFilterWriter routes logs to different files based on level
type levelFilterWriter struct {
	eventsFile *os.File
	errorsFile *os.File
}

func (w *levelFilterWriter) Write(p []byte) (n int, err error) {
	logLine := string(p)

	// Check if this is an error/warning level log
	isError := strings.Contains(logLine, " ERRO ") ||
		strings.Contains(logLine, " WARN ") ||
		strings.Contains(logLine, " FATA ")

	if isError {
		return w.errorsFile.Write(p)
	}
	return w.eventsFile.Write(p)
}

// init Logger initializes the logger singleton
func initLogger() {
	// Create logs directories if they don't exist
	if err := os.MkdirAll(filepath.Join("logs", "events"), 0755); err != nil {
		panic("Failed to create logs/events directory: " + err.Error())
	}
	if err := os.MkdirAll(filepath.Join("logs", "errors"), 0755); err != nil {
		panic("Failed to create logs/errors directory: " + err.Error())
	}

	// Open events log file
	var err error
	eventsFile, err = os.OpenFile(
		filepath.Join("logs", "events", "events.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic("Failed to open events log file: " + err.Error())
	}

	// Open errors log file
	errorsFile, err = os.OpenFile(
		filepath.Join("logs", "errors", "errors.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic("Failed to open errors log file: " + err.Error())
	}

	// Create logger with terminal output
	instance = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05",
		ReportCaller:    false,
	})

	// Create custom writer that routes logs by level
	levelWriter := &levelFilterWriter{
		eventsFile: eventsFile,
		errorsFile: errorsFile,
	}

	// Write to terminal + level-based files
	instance.SetOutput(io.MultiWriter(os.Stderr, levelWriter))
}

// Get returns the global logger instance
func Get() *log.Logger {
	once.Do(initLogger)
	return instance
}

// Close closes all log files - call this on graceful shutdown
func Close() {
	if eventsFile != nil {
		eventsFile.Close()
	}
	if errorsFile != nil {
		errorsFile.Close()
	}
}
