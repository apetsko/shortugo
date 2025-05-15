package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestShortenBatch_GRPC(t *testing.T) {
	userID := "user123"
	empty := ""
	one := "1"
	two := "2"
	example := "http://example.com"
	test := "http://test.com"
	badreq := "Bad Request: Empty URL"

	logger, _ := logging.New(zapcore.DebugLevel)

	tests := []struct {
		mockStorageSetup func(mockStorage *mocks.Storage)
		request          *pb.ShortenBatchRequest
		expectedBody     *pb.ShortenBatchResponse
		name             string
		userID           string
		expectedStatus   codes.Code
	}{
		{
			name:   "successful batch",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("PutBatch", mock.Anything, mock.Anything).Return(nil)
			},
			request: &pb.ShortenBatchRequest{
				UserId: &userID,
				Urls: []*pb.URLPair{
					{CorrelationId: &one, OriginalUrl: &example},
					{CorrelationId: &two, OriginalUrl: &test},
				},
			},
			expectedStatus: codes.OK,
			expectedBody:   nil, // проверим позже вручную
		},
		{
			name:   "empty url",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("PutBatch", mock.Anything, mock.Anything).Return(nil)
			},
			request: &pb.ShortenBatchRequest{
				UserId: &userID,
				Urls: []*pb.URLPair{
					{CorrelationId: &one, OriginalUrl: &empty},
				},
			},
			expectedStatus: codes.OK,
			expectedBody: &pb.ShortenBatchResponse{
				Results: []*pb.URLPair{
					{CorrelationId: &one, ShortUrl: &badreq},
				},
			},
		},
		{
			name:   "storage error",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("PutBatch", mock.Anything, mock.Anything).Return(errors.New("fail"))
			},
			request: &pb.ShortenBatchRequest{
				UserId: &userID,
				Urls: []*pb.URLPair{
					{CorrelationId: &one, OriginalUrl: &example},
				},
			},
			expectedStatus: codes.Internal,
			expectedBody:   nil,
		},
		{
			name:             "missing user id",
			userID:           "",
			mockStorageSetup: func(s *mocks.Storage) {},
			request: &pb.ShortenBatchRequest{
				UserId: &empty,
				Urls: []*pb.URLPair{
					{CorrelationId: &one, OriginalUrl: &example},
				},
			},
			expectedStatus: codes.InvalidArgument,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			if tt.mockStorageSetup != nil {
				tt.mockStorageSetup(mockStorage)
			}

			urlHandler := &httph.URLHandler{
				Storage: mockStorage,
				Logger:  logger,
				BaseURL: "http://short.ly",
			}

			grpcHandler := NewHandler(urlHandler)
			conn, cleanup, err := startGRPCServer(grpcHandler)
			require.NoError(t, err)
			defer cleanup()

			client := pb.NewURLShortenerClient(conn)

			resp, err := client.ShortenBatch(context.Background(), tt.request)

			if tt.expectedStatus == codes.OK {
				require.NoError(t, err)
				if tt.expectedBody != nil {
					assert.Len(t, resp.Results, len(tt.expectedBody.Results))
					for i := range resp.Results {
						assert.Equal(t, tt.expectedBody.Results[i].CorrelationId, resp.Results[i].CorrelationId)
						assert.Equal(t, tt.expectedBody.Results[i].ShortUrl, resp.Results[i].ShortUrl)
					}
				}
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedStatus, st.Code())
			}
		})
	}
}
