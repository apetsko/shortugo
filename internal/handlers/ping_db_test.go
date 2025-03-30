package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestPingDB(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(mockStorage *mocks.Storage)
		expectedStatus int
	}{
		{
			name: "database is reachable",
			mockSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Ping").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "database is unreachable",
			mockSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Ping").Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			logger, _ := logging.New(zapcore.DebugLevel)

			h := &URLHandler{
				storage: mockStorage,
				Logger:  logger,
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/ping", nil)

			h.PingDB(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
