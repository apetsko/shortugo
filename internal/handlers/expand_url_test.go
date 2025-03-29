package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

func BenchmarkExpandURL(b *testing.B) {
	logger, _ := logging.New(zapcore.DebugLevel)
	mockStorage := new(mocks.Storage)
	h := &URLHandler{
		storage: mockStorage,
		Logger:  logger,
	}

	testID := "abc123"
	mockURL := "https://example.com"

	mockStorage.On("Get", mock.Anything, testID).Return(mockURL, nil)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/"+testID, nil)
		w := httptest.NewRecorder()

		h.ExpandURL(w, req)

		resp := w.Result()

		assert.Equal(b, http.StatusTemporaryRedirect, resp.StatusCode, "unexpected status code")
		assert.Equal(b, mockURL, resp.Header.Get("Location"), "unexpected Location header")
	}
}

func TestURLHandler_ExpandURL(t *testing.T) {
	logger, err := logging.New(zapcore.DebugLevel)
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
