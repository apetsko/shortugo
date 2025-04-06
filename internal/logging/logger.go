package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogEntry defines the interface for log entries.
type LogEntry interface {
	// Write logs the status, bytes, header, elapsed time, and extra information.
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{})
	// Panic logs a panic message with the stack trace.
	Panic(v interface{}, stack []byte)
}

// Logger wraps the zap.SugaredLogger to provide structured logging.
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new Logger instance with the specified log level.
func New(level zapcore.Level) (*Logger, error) {
	// Configure the logger with development settings.
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level.SetLevel(level)

	// Build the logger with caller information.
	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Logger{logger.Sugar()}, nil
}

// Close syncs the logger, flushing any buffered log entries.
func (l *Logger) Close() error {
	return l.Sync()
}

// Debug logs a debug message with additional context.
func (l *Logger) Debug(message string, keysAndValues ...interface{}) {
	l.Debugw(message, keysAndValues...)
}

// Info logs an informational message with additional context.
func (l *Logger) Info(message string, keysAndValues ...interface{}) {
	l.Infow(message, keysAndValues...)
}

// Error logs an error message with additional context.
func (l *Logger) Error(message string, keysAndValues ...interface{}) {
	l.Errorw(message, keysAndValues...)
}

// Fatal logs a fatal message with additional context and then exits the application.
func (l *Logger) Fatal(message string, keysAndValues ...interface{}) {
	l.Fatalw(message, keysAndValues...)
}

// Printf logs a formatted informational message.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}
