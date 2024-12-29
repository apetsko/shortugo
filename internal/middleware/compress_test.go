package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decompressResponseBody(t *testing.T, body *bytes.Buffer) string {
	reader, err := gzip.NewReader(body)
	require.NoError(t, err, "Failed to create gzip reader")
	defer reader.Close()

	responseBody, err := io.ReadAll(reader)
	require.NoError(t, err, "Failed to read response body")

	return string(responseBody)
}

func TestGzipMiddleware_ResponseCompression_JSON(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, JSON World!"}`))
	})

	handler := GzipMiddleware()(testHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()

	var buf bytes.Buffer
	rec.Body = &buf

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"), "Expected Content-Encoding to be 'gzip'")

	responseBody := decompressResponseBody(t, &buf)
	assert.Equal(t, `{"message": "Hello, JSON World!"}`, responseBody, "Response body mismatch")
}

func TestGzipMiddleware_ResponseCompression_HTML(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<h1>Hello, HTML World!</h1>`))
	})

	handler := GzipMiddleware()(testHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()

	var buf bytes.Buffer
	rec.Body = &buf

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"), "Expected Content-Encoding to be 'gzip'")

	responseBody := decompressResponseBody(t, &buf)
	assert.Equal(t, `<h1>Hello, HTML World!</h1>`, responseBody, "Response body mismatch")
}
