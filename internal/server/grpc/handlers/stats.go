package handlers

import (
	"context"
	"net"
	"strings"

	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stats handles the gRPC request to retrieve internal usage statistics.
// Access is restricted to clients within a configured trusted subnet.
// It validates the caller's IP (provided in the request) and,
// if authorized, returns the total number of shortened URLs and unique users.
//
// Request:
//   - pb.StatsRequest with IP address in string format
//
// Response:
//   - pb.StatsResponse with UrlCount and UserCount
//
// Errors:
//   - PermissionDenied if TrustedSubnet is not configured, IP is invalid,
//     or the IP is outside the allowed subnet
//   - Internal if fetching stats from storage fails
func (h *Handler) Stats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	if h.URLHandler.TrustedSubnet == nil {
		h.URLHandler.Logger.Error("Forbidden. TrustedSubnet is not configured")
		return nil, status.Error(codes.PermissionDenied, "trusted subnet required")
	}

	ip := net.ParseIP(strings.TrimSpace(req.Ip))
	if ip == nil {
		h.URLHandler.Logger.Error("Forbidden: Invalid IP in request")
		return nil, status.Error(codes.PermissionDenied, "invalid IP")
	}

	if !h.URLHandler.TrustedSubnet.Contains(ip) {
		h.URLHandler.Logger.Errorf("Forbidden: IP %s not in trusted subnet", ip)
		return nil, status.Error(codes.PermissionDenied, "IP not allowed")
	}

	stats, err := h.URLHandler.Storage.Stats(ctx)
	if err != nil {
		h.URLHandler.Logger.Error("Failed to retrieve stats: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to retrieve stats")
	}

	return &pb.StatsResponse{
		UrlCount:  uint64(stats.Urls),
		UserCount: uint64(stats.Users),
	}, nil
}
