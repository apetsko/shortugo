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
		Auth:     mockAuth,
		Logger:   logger,
		Secret:   "some-Secret",
		ToDelete: make(chan models.BatchDeleteRequest, 100),
	}

	testUserID := "user-id"
	testIDs := []string{"abc123", "xyz789", "def456"}
	mockBody, _ := json.Marshal(testIDs)

	mockAuth.On("CookieGetUserID", mock.Anything, "some-Secret").Return(testUserID, nil)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/delete", bytes.NewReader(mockBody))
		w := httptest.NewRecorder()

		h.DeleteUserURLs(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(b, http.StatusAccepted, resp.StatusCode, "unexpected status code")

		body, _ := io.ReadAll(resp.Body)
		var responseIDs []string
		err := json.Unmarshal(body, &responseIDs)
		assert.NoError(b, err, "unexpected error while unmarshalling response")
		assert.Equal(b, testIDs, responseIDs, "unexpected IDs in response")
	}
}
func TestDeleteUserURLs(t *testing.T) {
	mockAuth := new(mocks.Authenticator)
	logger, _ := logging.New(zapcore.DebugLevel)
	h := &URLHandler{
		Auth:     mockAuth,
		Secret:   "valid",
		ToDelete: make(chan models.BatchDeleteRequest, 1),
		Logger:   logger,
	}

	tests := []struct {
		name           string
		setup          func()
		reqBody        io.Reader
		expectedStatus int
		validate       func(*httptest.ResponseRecorder)
	}{
		{
			name: "unauthorized user",
			setup: func() {
				mockAuth.On("CookieGetUserID", mock.Anything, "invalid").Return("", http.ErrNoCookie)
				h.Secret = "invalid"
			},
			reqBody:        http.NoBody,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid request body",
			setup: func() {
				mockAuth.On("CookieGetUserID", mock.Anything, "valid").Return("", nil)
				h.Secret = "valid"
			},
			reqBody:        bytes.NewReader([]byte("invalid-json")),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "successful deletion",
			setup: func() {
				mockAuth.On("CookieGetUserID", mock.Anything, "valid").Return("", nil)
				h.Secret = "valid"
			},
			reqBody: func() io.Reader {
				ids := []string{"id1", "id2"}
				body, _ := json.Marshal(ids)
				return bytes.NewReader(body)
			}(),
			expectedStatus: http.StatusAccepted,
			validate: func(w *httptest.ResponseRecorder) {
				var responseIDs []string
				body, _ := io.ReadAll(w.Body)
				err := json.Unmarshal(body, &responseIDs)
				assert.NoError(t, err)
				assert.Equal(t, []string{"id1", "id2"}, responseIDs)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/", tt.reqBody)
			h.DeleteUserURLs(w, r)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validate != nil {
				tt.validate(w)
			}
		})
	}
}
