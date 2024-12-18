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

func decompressGzip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TestGzipMiddleware_JSON(t *testing.T) {

	jsonHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, JSON!"}`))
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	handler := GzipMiddleware(jsonHandler)

	handler.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	require.Equal(t, res.Header.Get("Content-Encoding"), "gzip")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "Error reading response body")

	expected := `{"message": "Hello, JSON!"}`
	decompressed, err := decompressGzip(body)
	require.NoError(t, err, "Failed to decompress response")
	assert.Equal(t, expected, decompressed)

}

func TestGzipMiddleware_HTML(t *testing.T) {

	htmlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Hello, HTML!</h1></body></html>"))
	})

	req := httptest.NewRequest(http.MethodGet, "/html", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	handler := GzipMiddleware(htmlHandler)

	handler.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding to be gzip, got %s", res.Header.Get("Content-Encoding"))
	}

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "Error reading response body")

	expected := "<html><body><h1>Hello, HTML!</h1></body></html>"
	decompressed, err := decompressGzip(body)
	require.NoError(t, err, "Failed to decompress response")
	assert.Equal(t, expected, decompressed)
}

func TestGzipMiddleware_Text(t *testing.T) {

	textHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, plain text!"))
	})

	req := httptest.NewRequest(http.MethodGet, "/text", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	handler := GzipMiddleware(textHandler)

	handler.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.Header.Get("Content-Encoding") == "gzip" {
		t.Errorf("Expected no gzip encoding for text response")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello, plain text!"
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}
