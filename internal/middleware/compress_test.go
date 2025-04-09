package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func decompressGzip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decompressed), nil
}

func TestGzipMiddleware(t *testing.T) {
	testCases := []struct {
		name         string
		contentType  string
		body         string
		expectGzip   bool
		expectedBody string
	}{
		{
			name:         "JSON response",
			contentType:  "application/json",
			body:         `{"message": "Hello, JSON!"}`,
			expectGzip:   true,
			expectedBody: `{"message": "Hello, JSON!"}`,
		},
		{
			name:         "HTML response",
			contentType:  "text/html",
			body:         "<html><body><h1>Hello, HTML!</h1></body></html>",
			expectGzip:   true,
			expectedBody: "<html><body><h1>Hello, HTML!</h1></body></html>",
		},
		{
			name:         "Plain text response (no gzip)",
			contentType:  "text/plain",
			body:         "Hello, plain text!",
			expectGzip:   false,
			expectedBody: "Hello, plain text!",
		},
	}

	logger, err := logging.New(zapcore.DebugLevel)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				_, err := w.Write([]byte(tc.body))
				require.NoError(t, err, "error writing test response")
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Encoding", "gzip")

			rr := httptest.NewRecorder()
			middleware := GzipMiddleware(logger)
			middleware(handler).ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if tc.expectGzip {
				assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"), "Expected gzip encoding in response")
				compressedBody, err := io.ReadAll(res.Body)
				require.NoError(t, err, "Error reading response body")

				decompressed, err := decompressGzip(compressedBody)
				require.NoError(t, err, "Failed to decompress response")
				assert.Equal(t, tc.expectedBody, decompressed)
			} else {
				assert.Empty(t, res.Header.Get("Content-Encoding"), "Expected no gzip encoding for response")
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err, "Error reading response body")
				assert.Equal(t, tc.expectedBody, string(body))
			}
		})
	}
}
