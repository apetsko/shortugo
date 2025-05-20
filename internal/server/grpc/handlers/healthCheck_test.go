package handlers

import (
	"context"
	"testing"

	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_HealthCheck(t *testing.T) {
	h := &Handler{}

	req := &pb.HealthCheckRequest{}
	resp, err := h.HealthCheck(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Ok", resp.GetStatus())
}
