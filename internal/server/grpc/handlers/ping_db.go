package handlers

import (
	"context"

	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Ping checks the availability of the storage backend.
// It returns "OK" if the storage is reachable, otherwise returns a gRPC error.
func (h *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	if err := h.URLHandler.Storage.Ping(); err != nil {
		h.URLHandler.Logger.Error("Storage ping failed: " + err.Error())
		return nil, status.Error(codes.Unavailable, "storage unavailable")
	}
	return &pb.PingResponse{Status: "OK"}, nil
}
