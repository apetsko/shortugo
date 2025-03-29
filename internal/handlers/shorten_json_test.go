package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func BenchmarkShortenJSON(b *testing.B) {
	logger, _ := logging.New(zapcore.DebugLevel)
	mockStorage := new(mocks.Storage)
	mockAuth := new(mocks.Authenticator)
	h := &URLHandler{
		storage: mockStorage,
		Logger:  logger,
		baseURL: "http://localhost",
		secret:  "some-secret",
		auth:    mockAuth,
	}

	mockAuth.On("UserIDFromCookie", mock.Anything, "some-secret").Return("user-id", nil)
	mockStorage.On("Get", mock.Anything, mock.Anything).Return("", shared.ErrNotFound)
	mockStorage.On("Put", mock.Anything, mock.Anything).Return(nil)

	record := models.URLRecord{
		URL: "https://example.com",
	}
	requestBody, _ := json.Marshal(record)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		h.ShortenJSON(w, req)

		resp := w.Result()
		assert.Equal(b, http.StatusCreated, resp.StatusCode, "Status code should be 201 Created")

		body, err := io.ReadAll(resp.Body)
		assert.NoError(b, err, "Reading response body should not return an error")
		resp.Body.Close()

		var result models.Result
		err = json.Unmarshal(body, &result)
		assert.NoError(b, err, "Unmarshaling response body should not return an error")
		assert.Contains(b, result.Result, h.baseURL, "Result should contain the base URL")
		assert.Contains(b, result.Result, utils.GenerateID(record.URL, 8), "Result should contain the generated ID")
	}
}

func TestURLHandler_ShortenJSON(t *testing.T) {
	logger, err := logging.New(zapcore.DebugLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger, "fortytwo")
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
			record: models.URLRecord{URL: "https://practicum.yandex.ru/,", UserID: "55"},
			want: want{
				code:        201,
				ID:          "http://localhost:8080/3P9NwpqM",
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
