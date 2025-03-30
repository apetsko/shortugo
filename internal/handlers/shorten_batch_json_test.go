package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

	mockAuth.On("CookieGetUserID", mock.Anything, "some-secret").Return("user-id", nil)
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

func TestShortenBatchJSON(t *testing.T) {
	tests := []struct {
		name             string
		mockAuthSetup    func(mockAuth *mocks.Authenticator)
		mockStorageSetup func(mockStorage *mocks.Storage)
		requestBody      string
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "successful batch shortening",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("PutBatch", mock.Anything, mock.Anything).Return(nil)
			},
			requestBody: `[
				{"correlation_id":"1", "original_url":"http://example.com"},
				{"correlation_id":"2", "original_url":"http://test.com"}
			]`,
			expectedStatus: http.StatusCreated,
			expectedBody: `[
				{"correlation_id":"1", "short_url":"http://short.ly/"},
				{"correlation_id":"2", "short_url":"http://short.ly/"}
			]`,
		},
		{
			name: "bad request on invalid JSON",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      `invalid-json`,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name: "bad request on empty URL",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("PutBatch", mock.Anything, mock.Anything).Return(nil)
			},
			requestBody:    `[{"correlation_id":"1", "original_url":""}]`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `[{"correlation_id":"1", "short_url":"Bad Request: Empty URL"}]`,
		},
		{
			name: "internal server error on auth failure",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("", errors.New("auth error"))
				mockAuth.On("CookieSetUserID", mock.Anything, mock.Anything).Return("", errors.New("auth error"))
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      `[{"correlation_id":"1", "original_url":"http://example.com"}]`,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name: "storage failure",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("PutBatch", mock.Anything, mock.Anything).Return(errors.New("storage error"))
			},
			requestBody:    `[{"correlation_id":"1", "original_url":"http://example.com"}]`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mocks.Authenticator)
			mockStorage := new(mocks.Storage)
			logger, _ := logging.New(zapcore.DebugLevel)

			h := &URLHandler{
				auth:    mockAuth,
				storage: mockStorage,
				Logger:  logger,
				baseURL: "http://short.ly",
			}

			if tt.mockAuthSetup != nil {
				tt.mockAuthSetup(mockAuth)
			}
			if tt.mockStorageSetup != nil {
				tt.mockStorageSetup(mockStorage)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.requestBody))
			r.Header.Set("Content-Type", "application/json")

			h.ShortenBatchJSON(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				var expectedJSON, actualJSON []map[string]string
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedJSON)
				require.NoError(t, err)
				err = json.Unmarshal(w.Body.Bytes(), &actualJSON)
				require.NoError(t, err)

				for i := range expectedJSON {
					shortURL := actualJSON[i]["short_url"]

					if strings.HasPrefix(shortURL, "http://short.ly/") {
						expectedJSON[i]["short_url"] += shortURL[len("http://short.ly/"):]
					}
				}
			}
		})
	}
}
