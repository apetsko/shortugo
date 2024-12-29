package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/logging"
)

var zlogger, _ = logging.NewZapLogger()

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{ReadCloser: r, zr: zr}, nil
}

func (cr *compressReader) Read(p []byte) (int, error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
	return cr.zr.Close()
}

// bufferedResponseWriter — обертка для ResponseWriter, которая буферизует ответ
type bufferedResponseWriter struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func (brw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return brw.buffer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
				return
			}
			defer cr.Close()
			r.Body = cr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
		}

		buf := &bytes.Buffer{}
		bufferedWriter := &bufferedResponseWriter{ResponseWriter: w, buffer: buf}

		next.ServeHTTP(bufferedWriter, r)

		contentType := w.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			_, err := gz.Write(buf.Bytes())
			if err != nil {
				zlogger.Error(err.Error())
			}
		} else {
			_, err := w.Write(buf.Bytes())
			if err != nil {
				zlogger.Error(err.Error())
			}
		}
	})
}
