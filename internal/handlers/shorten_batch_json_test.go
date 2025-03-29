package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

func BenchmarkShortenBatchJSON(b *testing.B) {
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
	mockStorage.On("PutBatch", mock.Anything, mock.Anything).Return(nil)

	batchRequest := []models.BatchRequest{
		{ID: "1", OriginalURL: "https://example.com/page1"},
		{ID: "2", OriginalURL: "https://example.com/page2"},
		{ID: "3", OriginalURL: "https://example.com/page3"},
	}
	requestBody, _ := json.Marshal(batchRequest)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		h.ShortenBatchJSON(w, req)

		resp := w.Result()
		assert.Equal(b, http.StatusCreated, resp.StatusCode, "Status code should be 201 Created")

		body, err := io.ReadAll(resp.Body)
		assert.NoError(b, err, "Reading response body should not return an error")
		resp.Body.Close()

		var resps []models.BatchResponse
		err = json.Unmarshal(body, &resps)
		assert.NoError(b, err, "Unmarshaling response should not return an error")
		assert.Equal(b, len(batchRequest), len(resps), "Response length should match request length")

		for j, resp := range resps {
			assert.Equal(b, batchRequest[j].ID, resp.ID, "Response ID should match request ID")
			assert.Contains(b, resp.ShortURL, "http://localhost/", "ShortURL should contain base URL")
		}
	}
}
