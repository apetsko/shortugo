package server

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestHTTPS_Starts(t *testing.T) {
	go func() {
		storage := new(mocks.Storage)
		storage.On("Ping").Return(nil)
		logger, _ := logging.New(zapcore.DebugLevel)
		srv := &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			Addr:              ":8080",
			Handler:           Router(handlers.NewURLHandler("localhost", storage, logger, "secret", "")),
		}
		err := srv.ListenAndServeTLS("../../certs/cert.crt", "../../certs/cert.key")
		require.NoError(t, err)
	}()
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint
		},
	}

	resp, err := client.Get("https://localhost:8080/ping") // nolint
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
