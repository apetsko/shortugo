package http_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	pkghttp "github.com/apetsko/shortugo/internal/server/http"
	"github.com/apetsko/shortugo/internal/server/http/handlers"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestHTTPSServer_Run(t *testing.T) {
	logger, _ := logging.New(zapcore.DebugLevel)

	storage := new(mocks.Storage)
	storage.On("Ping").Return(nil)

	cfg := &config.Config{
		Host:        ":8443",
		EnableHTTPS: true,
		TLSCertPath: "../../../certs/cert.crt",
		TLSKeyPath:  "../../../certs/cert.key",
	}

	h := handlers.NewURLHandler("https://localhost:8443", storage, logger, "secret", "127.0.0.0/8")

	srv, err := pkghttp.Run(cfg, h, logger)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// TLS клиент с отключенной проверкой сертификата
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}

	var (
		resp   *http.Response
		reqErr error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://localhost:8443/ping", nil)
		resp, reqErr = client.Do(req)
		if reqErr == nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
		if resp != nil {
			if err = resp.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}
	}

	require.NoError(t, reqErr)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// Завершаем сервер
	require.NoError(t, srv.Shutdown(ctx))
}
