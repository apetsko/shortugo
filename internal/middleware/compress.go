// Package middleware provides HTTP middleware for the application.
// It includes functionality for request and response compression, logging, and other middleware utilities.
package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/apetsko/shortugo/internal/logging"
)

// sync.Pool for gzip.Reader
var gzipReaderPool = sync.Pool{
	New: func() any {
		return new(gzip.Reader)
	},
}

// bufferedResponseWriter buffers the response.
type bufferedResponseWriter struct {
	http.ResponseWriter
	buffer     *bytes.Buffer // Buffer to store the response body.
	statusCode int           // HTTP status code.
}

// WriteHeader sets the status code if it has not been set yet.
func (w *bufferedResponseWriter) WriteHeader(code int) {
	if w.statusCode == 0 {
		w.statusCode = code
	}
}

// Write sets the status code to http.StatusOK if it has not been set yet and writes the data to the buffer.
func (w *bufferedResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.buffer.Write(b)
}

// compressReader wraps the gzip.Reader to decompress the request body.
type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader // Gzip reader to decompress the data.
}

// newCompressReader creates a new compressReader by wrapping a gzip.Reader around the provided io.ReadCloser.
// It retrieves a gzip.Reader from the sync.Pool, resets it with the provided io.ReadCloser, and returns the compressReader.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr := gzipReaderPool.Get().(*gzip.Reader)
	err := zr.Reset(r)
	if err != nil {
		gzipReaderPool.Put(zr)
		return nil, err
	}
	return &compressReader{ReadCloser: r, zr: zr}, nil
}

// Read reads uncompressed data into p from the underlying gzip.Reader.
func (cr *compressReader) Read(p []byte) (int, error) {
	return cr.zr.Read(p)
}

// Close closes the underlying gzip.Reader and returns it to the sync.Pool.
func (cr *compressReader) Close() error {
	err := cr.zr.Close()
	gzipReaderPool.Put(cr.zr)
	return err
}

// GzipMiddleware is a middleware for gzip compression.
func GzipMiddleware(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Decompress the request body if needed
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
					return
				}
				defer func() {
					if err := cr.Close(); err != nil {
						logger.Error("Failed to close gzip reader: " + err.Error())
					}
				}()

				r.Body = cr
			}

			// If the client does not support gzip, just pass the request along
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Buffer the response
			buf := new(bytes.Buffer)
			bufferedWriter := &bufferedResponseWriter{
				ResponseWriter: w,
				buffer:         buf,
			}

			// Execute the request
			next.ServeHTTP(bufferedWriter, r)

			// Determine if compression is needed
			contentType := bufferedWriter.ResponseWriter.Header().Get("Content-Type")
			if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
				bufferedWriter.ResponseWriter.Header().Set("Content-Encoding", "gzip")
				bufferedWriter.ResponseWriter.WriteHeader(bufferedWriter.statusCode)

				gz := gzip.NewWriter(bufferedWriter.ResponseWriter)
				defer func() {
					if err := gz.Close(); err != nil {
						logger.Error("Failed to close gzip writer: " + err.Error())
					}
				}()

				if _, err := gz.Write(bufferedWriter.buffer.Bytes()); err != nil {
					logger.Error("Failed to compress response: " + err.Error())
					http.Error(bufferedWriter.ResponseWriter, "Failed to compress response", http.StatusInternalServerError)
					return
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			} else {
				// Send a regular response
				w.WriteHeader(bufferedWriter.statusCode)
				if _, err := w.Write(bufferedWriter.buffer.Bytes()); err != nil {
					logger.Error("Failed to write response: " + err.Error())
				}
			}
		})
	}
}
