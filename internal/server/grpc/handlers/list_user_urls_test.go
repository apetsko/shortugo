package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
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

func TestListUserURLs_GRPC(t *testing.T) {
	logger, _ := logging.New(zapcore.DebugLevel)

	tests := []struct {
		mockStorageSetup func(mockStorage *mocks.Storage)
		expectedBody     *pb.ListUserURLsResponse
		reqUserID        string
		name             string
		expectedStatus   codes.Code
	}{
		{
			name:      "no content",
			reqUserID: "user123",
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "user123").
					Return(nil, shared.ErrNotFound)
			},
			expectedStatus: codes.NotFound,
			expectedBody:   nil,
		},
		{
			name:      "internal error",
			reqUserID: "user123",
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "user123").
					Return(nil, errors.New("db error"))
			},
			expectedStatus: codes.Internal,
			expectedBody:   nil,
		},
		{
			name:      "successful retrieval",
			reqUserID: "user123",
			mockStorageSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("ListLinksByUserID", mock.Anything, "http://short.ly", "user123").
					Return([]models.URLRecord{
						{ID: "short1", URL: "http://example.com", UserID: "user123", Deleted: false},
						{ID: "short2", URL: "http://test.com", UserID: "user123", Deleted: false},
					}, nil)
			},
			expectedStatus: codes.OK,
			expectedBody: &pb.ListUserURLsResponse{
				Urls: []*pb.URLPair{
					{ShortUrl: "short1", OriginalUrl: "http://example.com"},
					{ShortUrl: "short2", OriginalUrl: "http://test.com"},
				},
			},
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
			ctx := context.Background()
			conn, cleanup, err := startGRPCServer(grpcHandler)
			require.NoError(t, err)
			defer cleanup()

			client := pb.NewURLShortenerClient(conn)

			resp, err := client.ListUserURLs(ctx, &pb.ListUserURLsRequest{
				UserId: tt.reqUserID,
			})

			if tt.expectedStatus == codes.OK {
				require.NoError(t, err)

				actualJSON, _ := json.Marshal(resp)
				expectedJSON, _ := json.Marshal(tt.expectedBody)

				require.JSONEq(t, string(expectedJSON), string(actualJSON))
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedStatus, st.Code())
			}
		})
	}
}
