package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

func BenchmarkListUserURLs(b *testing.B) {
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

	mockRecords := []models.URLRecord{
		{ID: "abc123", URL: "https://example.com/page1", UserID: "user-id", Deleted: false},
		{ID: "xyz789", URL: "https://example.com/page2", UserID: "user-id", Deleted: false},
		{ID: "def456", URL: "https://example.com/page3", UserID: "user-id", Deleted: false},
	}

	mockStorage.On("ListLinksByUserID", mock.Anything, "user-id", "http://localhost").Return(mockRecords, nil)
	mockAuth.On("CookieGetUserID", mock.Anything, "some-secret").Return("user-id", nil)

	cookie := &http.Cookie{Name: "shortugo", Value: "user-id"}

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/user/urls", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		h.ListUserURLs(w, req)

		if w.Result().StatusCode != http.StatusOK {
			b.Errorf("unexpected status code: got %v, want %v", w.Result().StatusCode, http.StatusOK)
		}

		expectedBody := `[{"short_url":"abc123","original_url":"https://example.com/page1"},{"short_url":"xyz789","original_url":"https://example.com/page2"},{"short_url":"def456","original_url":"https://example.com/page3"}]`
		if w.Body.String() != expectedBody {
			b.Errorf("unexpected body: got %v, want %v", w.Body.String(), expectedBody)
		}
	}
}

func TestListUserURLs(t *testing.T) {
	logger, _ := logging.New(zapcore.DebugLevel)

	tests := []struct {
		name             string
		mockAuthSetup    func(mockAuth *mocks.Authenticator)
		mockStorageSetup func(mockStorage *mocks.Storage)
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "successful retrieval",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "user123", "http://short.ly").
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
				mockStorage.On("ListLinksByUserID", mock.Anything, "user123", "http://short.ly").Return(nil, shared.ErrNotFound)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "internal server error on auth",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("", errors.New("auth error"))
				mockAuth.On("CookieSetUserID", mock.Anything, mock.Anything).Return("", errors.New("auth error"))
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {},
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name: "internal server error on storage",
			mockAuthSetup: func(mockAuth *mocks.Authenticator) {
				mockAuth.On("CookieGetUserID", mock.Anything, mock.Anything).Return("user123", nil)
			},
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "user123", "http://short.ly").
					Return(nil, errors.New("storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mocks.Authenticator)
			mockStorage := new(mocks.Storage)

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
