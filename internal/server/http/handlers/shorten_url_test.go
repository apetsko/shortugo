package handlers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func BenchmarkShortenURL(b *testing.B) {
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

	url := "https://example.com"
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader([]byte(url)))
		w := httptest.NewRecorder()

		h.ShortenURL(w, req)

		resp := w.Result()
		assert.Equal(b, http.StatusCreated, resp.StatusCode, "Status code should be 201 Created")

		body, err := io.ReadAll(resp.Body)
		assert.NoError(b, err, "Reading response body should not return an error")
		require.NoError(b, resp.Body.Close())

		shortenURL := string(body)
		assert.Contains(b, shortenURL, h.BaseURL, "Shorten URL should contain the base URL")
		assert.Contains(b, shortenURL, utils.GenerateID(url, 8), "Shorten URL should contain the generated ID")
	}
}

func TestURLHandler_ShortenURL(t *testing.T) {
	logger, err := logging.New(zapcore.DebugLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger, "fortytwo", "")
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

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.NoError(t, res.Body.Close())

			assert.Equal(t, test.want.ID, string(resBody))
		})
	}
}

func TestShortenURL(t *testing.T) {
	IDlen := 8
	shortenID := utils.GenerateID("http://example.com", IDlen)
	baseURL := "http://short.ly"
	shortenURL := baseURL + "/" + shortenID

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
			requestBody:    "http://example.com",
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
			requestBody:    "http://example.com",
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
			requestBody:      "http://example.com",
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name: "bad request on empty body",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      "",
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name: "bad request on empty URL",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			requestBody:      "",
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
			requestBody:    "http://example.com",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "error getting URL from Storage",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Get", mock.Anything, shortenID).Return("", errors.New("Storage error"))
			},
			requestBody:    "http://example.com",
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
			h.ShortenURL(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
