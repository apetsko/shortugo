package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestShortenJSON_GRPC(t *testing.T) {
	const (
		baseURL = "http://short.ly"
		url     = "http://example.com"
	)
	id := utils.GenerateID(url, 8)
	shortURL := baseURL + "/" + id

	tests := []struct {
		mockStorageSetup func(s *mocks.Storage)
		req              *pb.ShortenRequest
		name             string
		userID           string
		expectedShort    string
		expectedCode     codes.Code
	}{
		{
			name:   "successful URL shortening",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("Get", mock.Anything, id).Return("", shared.ErrNotFound)
				s.On("Put", mock.Anything, mock.Anything).Return(nil)
			},
			req: &pb.ShortenRequest{
				UserId:      "user123",
				OriginalUrl: url,
			},
			expectedCode:  codes.OK,
			expectedShort: shortURL,
		},
		{
			name:   "duplicate URL",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("Get", mock.Anything, id).Return(url, nil)
			},
			req: &pb.ShortenRequest{
				UserId:      "user123",
				OriginalUrl: url,
			},
			expectedCode:  codes.AlreadyExists,
			expectedShort: shortURL,
		},
		{
			name:             "missing user ID",
			userID:           "",
			mockStorageSetup: func(s *mocks.Storage) {},
			req: &pb.ShortenRequest{
				UserId:      "",
				OriginalUrl: url,
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:             "empty URL",
			userID:           "user123",
			mockStorageSetup: func(s *mocks.Storage) {},
			req: &pb.ShortenRequest{
				UserId:      "user123",
				OriginalUrl: "",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:   "storage error on get",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("Get", mock.Anything, id).Return("", errors.New("storage error"))
			},
			req: &pb.ShortenRequest{
				UserId:      "user123",
				OriginalUrl: url,
			},
			expectedCode: codes.Internal,
		},
		{
			name:   "storage error on put",
			userID: "user123",
			mockStorageSetup: func(s *mocks.Storage) {
				s.On("Get", mock.Anything, id).Return("", shared.ErrNotFound)
				s.On("Put", mock.Anything, mock.Anything).Return(errors.New("put failed"))
			},
			req: &pb.ShortenRequest{
				UserId:      "user123",
				OriginalUrl: url,
			},
			expectedCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			if tt.mockStorageSetup != nil {
				tt.mockStorageSetup(mockStorage)
			}

			logger, _ := logging.New(zapcore.DebugLevel)
			urlHandler := &httph.URLHandler{
				Storage: mockStorage,
				Logger:  logger,
				BaseURL: baseURL,
			}
			handler := NewHandler(urlHandler)
			conn, cleanup, err := startGRPCServer(handler)
			require.NoError(t, err)
			defer cleanup()

			client := pb.NewURLShortenerClient(conn)

			resp, err := client.ShortenJSON(context.Background(), tt.req)

			if tt.expectedCode == codes.OK {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedShort, resp.ShortUrl)
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())
			}
		})
	}
}
