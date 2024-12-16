package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

var compressibleTypes = map[string]bool{
	"application/json": true,
	"text/html":        true,
}

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	return cw.Writer.Write(b)
}

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
	if err := cr.zr.Close(); err != nil {
		return err
	}
	return cr.ReadCloser.Close()
}
func GzipMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				contentType := r.Header.Get("Content-Type")
				log.Println("Content-Type:", contentType)

				if compressibleTypes[contentType] {

					gz := gzip.NewWriter(w)
					defer gz.Close()

					w.Header().Set("Content-Encoding", "gzip")

					cw := &compressWriter{ResponseWriter: w, Writer: gz}

					next.ServeHTTP(cw, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
