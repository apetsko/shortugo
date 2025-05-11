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
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestPingDB_GRPC(t *testing.T) {
	tests := []struct {
		mockSetup    func(mockStorage *mocks.Storage)
		expectedBody *pb.PingResponse
		name         string
		expectedCode codes.Code
	}{
		{
			name: "database is reachable",
			mockSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Ping").Return(nil)
			},
			expectedCode: codes.OK,
			expectedBody: &pb.PingResponse{Status: "OK"},
		},
		{
			name: "database is unreachable",
			mockSetup: func(mockStorage *mocks.Storage) {
				mockStorage.On("Ping").Return(errors.New("storage error"))
			},
			expectedCode: codes.Unavailable,
			expectedBody: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			logger, _ := logging.New(zapcore.DebugLevel)

			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage)
			}

			urlHandler := &httph.URLHandler{
				Storage: mockStorage,
				Logger:  logger,
			}

			grpcHandler := NewHandler(urlHandler)
			conn, cleanup, err := startGRPCServer(grpcHandler)
			require.NoError(t, err)
			defer cleanup()

			client := pb.NewURLShortenerClient(conn)

			resp, err := client.Ping(context.Background(), &pb.PingRequest{})

			if tt.expectedCode == codes.OK {
				require.NoError(t, err)
				assert.True(t, proto.Equal(tt.expectedBody, resp))
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())
			}
		})
	}
}
