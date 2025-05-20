package handlers

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/models"
	httph "github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestStats_GRPC(t *testing.T) {
	_, trustedNet, _ := net.ParseCIDR("192.168.0.0/24")
	ten := int64(10)
	five := int64(5)
	tests := []struct {
		trustedSubnet    *net.IPNet
		mockStats        *models.Stats
		expectedResponse *pb.StatsResponse
		mockError        error
		name             string
		ip               string
		expectedCode     codes.Code
	}{
		{
			name:          "success",
			ip:            "192.168.0.42",
			trustedSubnet: trustedNet,
			mockStats:     &models.Stats{Urls: 10, Users: 5},
			expectedCode:  codes.OK,
			expectedResponse: &pb.StatsResponse{
				UrlCount:  &ten,
				UserCount: &five,
			},
		},
		{
			name:          "no subnet configured",
			ip:            "192.168.0.42",
			trustedSubnet: nil,
			expectedCode:  codes.PermissionDenied,
		},
		{
			name:          "invalid IP",
			ip:            "invalid-ip",
			trustedSubnet: trustedNet,
			expectedCode:  codes.PermissionDenied,
		},
		{
			name:          "IP outside subnet",
			ip:            "10.0.0.1",
			trustedSubnet: trustedNet,
			expectedCode:  codes.PermissionDenied,
		},
		{
			name:          "storage returns error",
			ip:            "192.168.0.42",
			trustedSubnet: trustedNet,
			mockError:     errors.New("storage error"),
			expectedCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			logger, _ := logging.New(zapcore.DebugLevel)

			if tt.mockStats != nil || tt.mockError != nil {
				mockStorage.On("Stats", mock.Anything).Return(tt.mockStats, tt.mockError)
			}

			h := &httph.URLHandler{
				Storage:       mockStorage,
				Logger:        logger,
				TrustedSubnet: tt.trustedSubnet,
			}

			grpcHandler := NewHandler(h)
			conn, cleanup, err := startGRPCServer(grpcHandler)
			require.NoError(t, err)
			defer cleanup()

			client := pb.NewURLShortenerClient(conn)

			resp, err := client.Stats(context.Background(), &pb.StatsRequest{Ip: &tt.ip})

			if tt.expectedCode == codes.OK {
				require.NoError(t, err)
				assert.True(t, proto.Equal(tt.expectedResponse, resp))
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())
			}
		})
	}
}
