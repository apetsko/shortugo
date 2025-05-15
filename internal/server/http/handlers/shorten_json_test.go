package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
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
		Storage: mockStorage,
		Logger:  logger,
		BaseURL: "http://localhost",
		Secret:  "some-Secret",
		Auth:    mockAuth,
	}

	mockAuth.On("CookieGetUserID", mock.Anything, "some-Secret").Return("user-id", nil)
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
		require.NoError(b, resp.Body.Close())

		var result models.Result
		err = json.Unmarshal(body, &result)
		assert.NoError(b, err, "Unmarshaling response body should not return an error")
		assert.Contains(b, result.Result, h.BaseURL, "Result should contain the base URL")
		assert.Contains(b, result.Result, utils.GenerateID(record.URL, 8), "Result should contain the generated ID")
	}
}
func TestShortenJSON(t *testing.T) {
	IDlen := 8
	shortenID := utils.GenerateID("http://example.com", IDlen)
	baseURL := "http://short.ly"
	shortenURL := fmt.Sprintf(`{"result":"%s/%s"}`, baseURL, shortenID)

	tests := []struct {
		mockAuthSetup    func(mockAuth *mocks.Authenticator)
		mockStorageSetup func(mockStorage *mocks.Storage)
		name             string
		requestBody      string
		expectedBody     string
		expectedStatus   int
	}{
		{
			name: "successful URL shortening",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Get", mock.Anything, shortenID).Return("", shared.ErrNotFound)
				mockStorage.On("Put", mock.Anything, mock.Anything).Return(nil)
			},
			requestBody:    `{"url":"http://example.com"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   shortenURL,
		},
		{
			name: "duplicate URL returns conflict",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Get", mock.Anything, shortenID).Return("http://short.ly/shortID", nil)
			},
			requestBody:    `{"url":"http://example.com"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   shortenURL,
		},
		{
			name: "Auth error",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("", errors.New("Auth error"))
				mockAuth.On("CookieSetUserID", mock.Anything, mock.Anything).Return("", errors.New("Auth error"))
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      `{"url":"http://example.com"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name: "bad request on empty body",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      ``,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name: "bad request on invalid JSON",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      `{"url":123}`,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name: "bad request on empty URL",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      `{"url":""}`,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name: "error storing URL",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Get", mock.Anything, shortenID).Return("", shared.ErrNotFound)
				mockStorage.On("Put", mock.Anything, mock.Anything).Return(errors.New("Storage error"))
			},
			requestBody:    `{"url":"http://example.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error getting URL from Storage",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Get", mock.Anything, shortenID).Return("", errors.New("Storage error"))
			},
			requestBody:    `{"url":"http://example.com"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mocks.Authenticator)
			mockStorage := new(mocks.Storage)
			logger, _ := logging.New(zapcore.DebugLevel)

			h := &URLHandler{
				Auth:    mockAuth,
				Storage: mockStorage,
				Logger:  logger,
				BaseURL: baseURL,
			}

			if tt.mockAuthSetup != nil {
				tt.mockAuthSetup(mockAuth)
			}
			if tt.mockStorageSetup != nil {
				tt.mockStorageSetup(mockStorage)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			h.ShortenJSON(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				var expectedJSON, actualJSON map[string]string
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedJSON)
				assert.NoError(t, err)
				err = json.Unmarshal(w.Body.Bytes(), &actualJSON)
				assert.NoError(t, err)

				assert.Equal(t, expectedJSON, actualJSON)
			}
		})
	}
}
