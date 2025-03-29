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

func BenchmarkDeleteUserURLs(b *testing.B) {
	logger, _ := logging.New(zapcore.DebugLevel)
	mockAuth := new(mocks.Authenticator)
	h := &URLHandler{
		auth:     mockAuth,
		Logger:   logger,
		secret:   "some-secret",
		ToDelete: make(chan models.BatchDeleteRequest, 100), // канал для удаления
	}

	testUserID := "user-id"
	testIDs := []string{"abc123", "xyz789", "def456"}
	mockBody, _ := json.Marshal(testIDs)

	mockAuth.On("UserIDFromCookie", mock.Anything, "some-secret").Return(testUserID, nil)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/delete", bytes.NewReader(mockBody))
		w := httptest.NewRecorder()

		h.DeleteUserURLs(w, req)

		resp := w.Result()

		assert.Equal(b, http.StatusAccepted, resp.StatusCode, "unexpected status code")

		body, _ := io.ReadAll(resp.Body)
		var responseIDs []string
		err := json.Unmarshal(body, &responseIDs)
		assert.NoError(b, err, "unexpected error while unmarshalling response")
		assert.Equal(b, testIDs, responseIDs, "unexpected IDs in response")
	}
}
