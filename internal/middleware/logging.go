package middleware

import (
	"net/http"
	"time"
)

type logger interface {
	Info(message string, keysAndValues ...interface{})
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	logger logger
	http.ResponseWriter
	responseData *responseData
}

func newLoggingResponseWriter(w http.ResponseWriter, logger logger) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		logger:         logger,
		responseData:   &responseData{},
	}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(logger logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			uri := r.RequestURI
			method := r.Method

			lw := newLoggingResponseWriter(w, logger)

			next.ServeHTTP(lw, r)
			duration := time.Since(start)

			lw.logger.Info(
				"middleware logger",
				"uri", uri,
				"method", method,
				"status", lw.responseData.status,
				"duration", duration.Nanoseconds(),
				"size", lw.responseData.size,
			)
		}

		return http.HandlerFunc(logFn)
	}
}