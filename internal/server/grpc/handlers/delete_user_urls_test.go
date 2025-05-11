package handlers

import (
	"context"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDeleteUserURLs_GRPC(t *testing.T) {
	toDelete := make(chan models.BatchDeleteRequest, 1)
	logger, _ := logging.New(zapcore.DebugLevel)
	mockAuth := new(mocks.Authenticator)
	urlHandler := &httph.URLHandler{
		Storage:  nil, // not needed here
		ToDelete: toDelete,
		Secret:   "valid",
		Logger:   logger,
		Auth:     mockAuth,
	}

	grpcHandler := NewHandler(urlHandler)
	ctx := context.Background()
	conn, cleanup, err := startGRPCServer(grpcHandler)
	require.NoError(t, err)
	defer cleanup()

	client := pb.NewURLShortenerClient(conn)

	req := &pb.DeleteUserURLsRequest{
		UserId:      "test-user",
		ShortUrlIds: []string{"id1", "id2"},
	}

	resp, err := client.DeleteUserURLs(ctx, req)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	select {
	case msg := <-toDelete:
		assert.Equal(t, "test-user", msg.UserID)
		assert.Equal(t, []string{"id1", "id2"}, msg.Ids)
	default:
		t.Fatal("expected message on ToDelete channel")
	}
}
