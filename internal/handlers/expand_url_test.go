package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/stretchr/testify/assert"
)

func TestURLHandler_ExpandURL(t *testing.T) {
	logger, err := logging.New()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger, "fortytwo")
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
			record: models.URLRecord{ID: "QrPnX5IU", URL: "https://practicum.yandex.ru/", UserID: "55"},
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
