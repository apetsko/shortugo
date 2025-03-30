package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func TestLogMiddleware(t *testing.T) {
	mockLogger := new(MockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, "GET", args.Get(0).(string)) // Метод запроса
		fields := args.Get(1).([]interface{})

		assert.Contains(t, fields, "uri")
		assert.Contains(t, fields, "/test")
		assert.Contains(t, fields, "status")
		assert.Contains(t, fields, 200)
	}).Once()

	middleware := LogMiddleware(mockLogger)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("UserID", "12345")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())

	mockLogger.AssertExpectations(t)
}
