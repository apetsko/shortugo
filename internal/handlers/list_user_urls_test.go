package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
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
	mockAuth.On("UserIDFromCookie", mock.Anything, "some-secret").Return("user-id", nil)

	cookie := &http.Cookie{Name: "shortugo", Value: "user-id"} // Имя куки должно совпадать с функцией

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
