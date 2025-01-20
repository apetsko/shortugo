package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
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
		URL      string
	}
	tests := []struct {
		name   string
		record models.URLRecord
		want   want
	}{
		{
			name:   "positive test #1",
			record: models.URLRecord{ID: "QrPnX5IU", URL: "https://practicum.yandex.ru/"},
			want: want{
				code:     307,
				Location: "https://practicum.yandex.ru/",
				URL:      "http://localhost:8080/QrPnX5IU",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = handler.storage.Put(context.Background(), test.record)
			request := httptest.NewRequest(http.MethodGet, test.want.URL, nil)
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
		name   string
		record models.URLRecord
		want   want
	}{
		{
			name:   "positive test #1",
			record: models.URLRecord{URL: "https://practicum.yandex.ru/"},
			want: want{
				code:        201,
				ID:          "http://localhost:8080/QrPnX5IU",
				ContentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, err := json.Marshal(test.record)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, test.want.ID, bytes.NewBuffer(url))
			w := httptest.NewRecorder()
			handler.ShortenJSON(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			var resp models.Result
			err = json.Unmarshal(b, &resp)
			require.NoError(t, err)

			assert.Equal(t, test.want.ContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.ID, resp.Result)
		})
	}
}
