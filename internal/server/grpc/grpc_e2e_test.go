package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/server/http/handlers"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	pb "github.com/apetsko/shortugo/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPC_E2E(t *testing.T) {
	testIP := "127.0.0.1"
	testUser := "test-user"
	one := "1"

	mockStorage := new(mocks.Storage)
	logger, _ := logging.New(zapcore.DebugLevel)
	toDelete := make(chan models.BatchDeleteRequest, 1)

	// Setup common mocks
	mockStorage.On("Ping").Return(nil)
	mockStorage.On("PutBatch", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("DeleteUserURLs", mock.Anything, []string{"abc123"}, "test-user").Return(nil)
	mockStorage.On("Stats", mock.Anything).Return(&models.Stats{Urls: 2, Users: 1}, nil)
	mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "test-user").Return([]models.URLRecord{
		{ID: "abc123", URL: "http://example.com", UserID: "test-user"},
	}, nil)

	// gRPC-сервер
	cfg := &config.Config{
		GRPCHost:    "127.0.0.1:19091",
		EnableHTTPS: false,
	}
	handler := handlers.NewURLHandler("http://short.ly", mockStorage, logger, "secret", "127.0.0.0/8")
	handler.ToDelete = toDelete

	srv, err := Run(cfg, handler, logger)
	require.NoError(t, err)
	defer srv.GracefulStop()
	time.Sleep(200 * time.Millisecond)

	conn, err := grpc.NewClient("passthrough:///127.0.0.1:19091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	client := pb.NewURLShortenerClient(conn)
	ctx := context.Background()

	// Step 1: Ping
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	require.NoError(t, err)
	assert.Equal(t, "OK", resp.GetStatus())

	// Step 2: Shorten
	uniqueURL := fmt.Sprintf("http://example.com/e2e-%d", time.Now().UnixNano())
	shortenedID := utils.GenerateID(uniqueURL, 8)

	// Get for Shorten: simulate that ID not found
	mockStorage.On("Get", mock.Anything, shortenedID).Return("", shared.ErrNotFound).Once()
	mockStorage.On("Put", mock.Anything, mock.Anything).Return(nil).Once()

	shortenResp, err := client.Shorten(ctx, &pb.ShortenRequest{
		OriginalUrl: &uniqueURL,
		UserId:      &testUser,
	})
	require.NoError(t, err)
	assert.Equal(t, "http://short.ly/"+shortenedID, shortenResp.GetShortUrl())

	// Step 3: Expand — Get returns actual URL
	mockStorage.On("Get", mock.Anything, shortenedID).Return(uniqueURL, nil).Once()

	expandResp, err := client.Expand(ctx, &pb.ExpandRequest{
		ShortUrlId: &shortenedID,
	})
	require.NoError(t, err)
	assert.Equal(t, uniqueURL, expandResp.GetOriginalUrl())

	// Step 4: ShortenBatch
	batchURL := fmt.Sprintf("http://example.com/batch-%d", time.Now().UnixNano())

	batchResp, err := client.ShortenBatch(ctx, &pb.ShortenBatchRequest{
		UserId: &testUser,
		Urls: []*pb.URLPair{
			{CorrelationId: &one, OriginalUrl: &batchURL},
		},
	})
	require.NoError(t, err)
	require.Len(t, batchResp.GetResults(), 1)
	assert.Equal(t, "1", batchResp.GetResults()[0].GetCorrelationId())
	assert.Contains(t, batchResp.GetResults()[0].GetShortUrl(), "http://short.ly/")

	// Step 5: ListUserURLs
	listResp, err := client.ListUserURLs(ctx, &pb.ListUserURLsRequest{
		UserId: &testUser,
	})
	require.NoError(t, err)
	require.Len(t, listResp.Urls, 1)
	assert.Equal(t, "abc123", listResp.GetUrls()[0].GetShortUrl())

	// Step 6: DeleteUserURLs
	delResp, err := client.DeleteUserURLs(ctx, &pb.DeleteUserURLsRequest{
		UserId:      &testUser,
		ShortUrlIds: []string{"abc123"},
	})
	require.NoError(t, err)
	assert.True(t, delResp.GetSuccess())

	go func() {
		select {
		case req := <-toDelete:
			_ = mockStorage.DeleteUserURLs(context.Background(), req.Ids, req.UserID)
		case <-time.After(time.Second):
			t.Error("timeout waiting for delete request")
		}
	}()

	// Step 7: Stats
	statsResp, err := client.Stats(ctx, &pb.StatsRequest{Ip: &testIP})
	require.NoError(t, err)
	assert.Equal(t, int64(2), statsResp.GetUrlCount())
	assert.Equal(t, int64(1), statsResp.GetUserCount())

	mockStorage.AssertExpectations(t)
}
