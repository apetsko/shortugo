// Package middleware is middleware package
package middleware

import (
	"net/http"
	"time"
)

// logger defines the interface for logging.
type logger interface {
	// Info logs an informational message with additional context.
	Info(message string, keysAndValues ...interface{})
}

// responseData holds the HTTP response status and size.
type responseData struct {
	status int // HTTP status code.
	size   int // Size of the response body.
}

// logResponseWriter wraps the http.ResponseWriter to capture response data.
type logResponseWriter struct {
	http.ResponseWriter
	responseData *responseData // Captured response data.
}

// newLogResponseWriter creates a new logResponseWriter.
func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{
		ResponseWriter: w,
		responseData:   &responseData{status: http.StatusOK},
	}
}

// Write writes the data to the response and captures the size.
func (r *logResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader sets the HTTP status code and captures it.
func (r *logResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LogMiddleware logs the details of each HTTP request and response.
func LogMiddleware(logger logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture response data.
			lw := newLogResponseWriter(w)
			next.ServeHTTP(lw, r)

			// Log the request and response details.
			logger.Info(
				r.Method,
				"uri", r.RequestURI,
				"UserID", r.Header.Get("UserID"),
				"duration_ms", time.Since(start).Milliseconds(),
				"size", lw.responseData.size,
				"status", lw.responseData.status,
				"IP", r.RemoteAddr,
			)
		})
	}
}
