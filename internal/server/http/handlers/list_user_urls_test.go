package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func BenchmarkListUserURLs(b *testing.B) {
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

	mockRecords := []models.URLRecord{
		{ID: "abc123", URL: "https://example.com/page1", UserID: "user-id", Deleted: false},
		{ID: "xyz789", URL: "https://example.com/page2", UserID: "user-id", Deleted: false},
		{ID: "def456", URL: "https://example.com/page3", UserID: "user-id", Deleted: false},
	}

	mockStorage.On("ListLinksByUserID", mock.Anything, "http://localhost", "user-id").Return(mockRecords, nil)
	mockAuth.On("CookieGetUserID", mock.Anything, "some-Secret").Return("user-id", nil)

	cookie := &http.Cookie{Name: "shortugo", Value: "user-id"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/user/urls", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		h.ListUserURLs(w, req)

		resp := w.Result()
		defer require.NoError(b, resp.Body.Close())
		body, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		require.Equal(b, http.StatusOK, resp.StatusCode, "unexpected status code")

		expectedBody := `[{"short_url":"abc123","original_url":"https://example.com/page1"},{"short_url":"xyz789","original_url":"https://example.com/page2"},{"short_url":"def456","original_url":"https://example.com/page3"}]`
		require.JSONEq(b, expectedBody, string(body), "unexpected body")
	}
}

func TestListUserURLs(t *testing.T) {
	logger, _ := logging.New(zapcore.DebugLevel)

	tests := []struct {
		mockAuthSetup    func(mockAuth *mocks.Authenticator)
		mockStorageSetup func(mockStorage *mocks.Storage)
		name             string
		expectedBody     string
		expectedStatus   int
	}{
		{
			name: "successful retrieval",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "user123").
					Return([]models.URLRecord{
						{ID: "short1", URL: "http://example.com", UserID: "user123", Deleted: false},
						{ID: "short2", URL: "http://test.com", UserID: "user123", Deleted: false},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"short_url":"short1","original_url":"http://example.com"},{"short_url":"short2","original_url":"http://test.com"}]`,
		},
		{
			name: "no content",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "user123").Return(nil, shared.ErrNotFound)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "internal server error on Auth",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("", errors.New("Auth error"))
				mockAuth.On("CookieSetUserID", mock.Anything, mock.Anything).Return("", errors.New("Auth error"))
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name: "internal server error on Storage",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "user123").
					Return(nil, errors.New("Storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mocks.Authenticator)
			mockStorage := new(mocks.Storage)

			h := &URLHandler{
				Auth:    mockAuth,
				Storage: mockStorage,
				Logger:  logger,
				BaseURL: "http://short.ly",
			}

			if tt.mockAuthSetup != nil {
				tt.mockAuthSetup(mockAuth)
			}
			if tt.mockStorageSetup != nil {
				tt.mockStorageSetup(mockStorage)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			h.ListUserURLs(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			mockAuth.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}
