package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storage/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var logger = initLogger()

func initLogger() *logging.ZapLogger {
	logger, err := logging.NewZapLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	return logger
}
func TestURLHandler_ShortenURL(t *testing.T) {
	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger)
	type want struct {
		ID   string
		code int
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "positive test #1",
			URL:  "https://practicum.yandex.ru/",
			want: want{

				code: 201,
				ID:   "http://localhost:8080/QrPnX5IU",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.want.ID, strings.NewReader(test.URL))
			w := httptest.NewRecorder()
			handler.ShortenURL(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.ID, string(resBody))
		})
	}
}

func TestURLHandler_ExpandURL(t *testing.T) {
	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger)
	type want struct {
		code     int
		Location string
	}
	tests := []struct {
		shortenURL string
		id         string
		name       string
		want       want
	}{
		{
			name:       "positive test #1",
			shortenURL: "http://localhost:8080/QrPnX5IU",
			id:         "QrPnX5IU",
			want: want{
				code:     307,
				Location: "https://practicum.yandex.ru/",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = handler.storage.Put(test.id, test.want.Location)
			request := httptest.NewRequest(http.MethodGet, test.shortenURL, nil)
			w := httptest.NewRecorder()
			handler.ExpandURL(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.Location, res.Header.Get("Location"))
		})
	}
}

func TestURLHandler_ShortenJSON(t *testing.T) {
	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger)
	type want struct {
		ID          string
		code        int
		ContentType string
	}
	tests := []struct {
		name string
		URL  models.Request
		want want
	}{
		{
			name: "positive test #1",
			URL:  models.Request{URL: "https://practicum.yandex.ru/"},
			want: want{
				code:        201,
				ID:          "http://localhost:8080/QrPnX5IU",
				ContentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urljson, err := json.Marshal(test.URL)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, test.want.ID, bytes.NewBuffer(urljson))
			w := httptest.NewRecorder()
			handler.ShortenJSON(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			var resp models.Response
			err = json.Unmarshal(b, &resp)
			require.NoError(t, err)

			assert.Equal(t, test.want.ContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.ID, resp.Result)
		})
	}
}
