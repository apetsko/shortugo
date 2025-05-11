package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	"github.com/apetsko/shortugo/internal/storages/shared"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestExpand_GRPC(t *testing.T) {
	mockStorage := new(mocks.Storage)
	logger, _ := logging.New(zapcore.DebugLevel)

	urlHandler := &httph.URLHandler{
		Storage: mockStorage,
		Logger:  logger,
	}

	grpcHandler := NewHandler(urlHandler)
	ctx := context.Background()
	conn, cleanup, err := startGRPCServer(grpcHandler)
	require.NoError(t, err)
	defer cleanup()

	client := pb.NewURLShortenerClient(conn)

	tests := []struct {
		assertErr  func(*testing.T, error)
		mockError  error
		name       string
		id         string
		mockReturn string
	}{
		{
			name:       "successful resolve",
			id:         "abc123",
			mockReturn: "http://example.com",
			mockError:  nil,
			assertErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:       "url gone",
			id:         "expired123",
			mockReturn: "",
			mockError:  shared.ErrGone,
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, s.Code())
			},
		},
		{
			name:       "url not found",
			id:         "missing123",
			mockReturn: "",
			mockError:  shared.ErrNotFound,
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.NotFound, s.Code())
			},
		},
		{
			name:       "internal error",
			id:         "error123",
			mockReturn: "",
			mockError:  errors.New("db fail"),
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				s, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.On("Get", mock.Anything, tt.id).Return(tt.mockReturn, tt.mockError)

			resp, err := client.Expand(ctx, &pb.ExpandRequest{ShortUrlId: tt.id})

			if tt.mockError == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, resp.OriginalUrl)
			} else {
				require.Nil(t, resp)
				tt.assertErr(t, err)
			}
		})
	}
}
