package handlers

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestURLHandler_Stats(t *testing.T) {
	type mocksSetup struct {
		storageStats *models.Stats
		storageErr   error
	}

	_, trustedNet, _ := net.ParseCIDR("192.168.0.0/24")

	tests := []struct {
		setupMocks     mocksSetup
		trustedSubnet  *net.IPNet
		name           string
		ipHeader       string
		expectedBody   string
		expectedHeader string
		expectedCode   int
	}{
		{
			name:          "success",
			trustedSubnet: trustedNet,
			ipHeader:      "192.168.0.42",
			setupMocks: mocksSetup{
				storageStats: &models.Stats{Urls: 10, Users: 5},
			},
			expectedCode:   http.StatusOK,
			expectedBody:   `{"urls":10,"users":5}`,
			expectedHeader: "application/json",
		},
		{
			name:          "nil subnet",
			trustedSubnet: nil,
			ipHeader:      "192.168.0.42",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "invalid IP",
			trustedSubnet: trustedNet,
			ipHeader:      "invalid-ip",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "IP not in subnet",
			trustedSubnet: trustedNet,
			ipHeader:      "10.0.0.1",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "storage returns error",
			trustedSubnet: trustedNet,
			ipHeader:      "192.168.0.42",
			setupMocks: mocksSetup{
				storageErr: errors.New("db down"),
			},
			expectedCode:   http.StatusOK,
			expectedBody:   `null`,
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			logger, _ := logging.New(zapcore.DebugLevel)

			if tt.setupMocks.storageStats != nil || tt.setupMocks.storageErr != nil {
				mockStorage.On("Stats", mock.Anything).Return(tt.setupMocks.storageStats, tt.setupMocks.storageErr)
			}

			h := &URLHandler{
				Storage:       mockStorage,
				Logger:        logger,
				TrustedSubnet: tt.trustedSubnet,
			}

			req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			req.Header.Set("X-Real-IP", tt.ipHeader)
			w := httptest.NewRecorder()

			h.Stats(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				require.NoError(t, err)
			}()

			assert.Equal(t, tt.expectedCode, res.StatusCode)

			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, res.Header.Get("Content-Type"))
			}

			if tt.expectedBody != "" {
				bodyBytes, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				body := strings.TrimSpace(string(bodyBytes))
				assert.JSONEq(t, tt.expectedBody, body)
			}
		})
	}
}
