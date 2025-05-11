package handlers

import (
	"context"

	pb "github.com/apetsko/shortugo/proto"
)

// HealthCheck handles a health check request.
// It responds with a static "OK" status to confirm the service is alive.
//
// This method corresponds to a simple HTTP healthcheck endpoint.
func (h *Handler) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: "OK"}, nil
}
