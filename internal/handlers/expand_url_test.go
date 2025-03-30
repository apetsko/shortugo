package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/storages/shared"
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

func TestExpandURL(t *testing.T) {
	mockStorage := new(mocks.Storage)
	logger, _ := logging.New(zapcore.DebugLevel)
	h := &URLHandler{
		storage: mockStorage,
		Logger:  logger,
	}

	tests := []struct {
		name           string
		urlID          string
		mockReturn     string
		mockError      error
		expectedStatus int
		validate       func(*httptest.ResponseRecorder)
	}{
		{
			name:           "successful redirect",
			urlID:          "test-id",
			mockReturn:     "http://example.com",
			mockError:      nil,
			expectedStatus: http.StatusTemporaryRedirect,
			validate: func(w *httptest.ResponseRecorder) {
				assert.Equal(t, "http://example.com", w.Header().Get("Location"))
			},
		},
		{
			name:           "URL gone",
			urlID:          "expired-id",
			mockReturn:     "",
			mockError:      shared.ErrGone,
			expectedStatus: http.StatusGone,
		},
		{
			name:           "URL not found",
			urlID:          "missing-id",
			mockReturn:     "",
			mockError:      shared.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal error",
			urlID:          "error-id",
			mockReturn:     "",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.On("Get", mock.Anything, tt.urlID).Return(tt.mockReturn, tt.mockError)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/"+tt.urlID, nil)

			h.ExpandURL(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.validate != nil {
				tt.validate(w)
			}
		})
	}
}
