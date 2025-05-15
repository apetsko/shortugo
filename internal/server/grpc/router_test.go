package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouterGRPC(t *testing.T) {
	// Можно передать nil, потому что NewHandler допускает nil
	srv := RouterGRPC(nil)

	require.NotNil(t, srv)
}
