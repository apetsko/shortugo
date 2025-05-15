package handlers

import (
	"context"
	"errors"

	"github.com/apetsko/shortugo/internal/storages/shared"
	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Expand resolves a short URL ID to its original URL.
// Returns gRPC status codes based on the error encountered.
func (h *Handler) Expand(ctx context.Context, req *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	originalURL, err := h.URLHandler.Storage.Get(ctx, req.GetShortUrlId())
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrGone):
			h.URLHandler.Logger.Error("URL is gone: " + req.GetShortUrlId())
			return nil, status.Error(codes.FailedPrecondition, "URL is gone")

		case errors.Is(err, shared.ErrNotFound):
			h.URLHandler.Logger.Error("URL not found: " + req.GetShortUrlId())
			return nil, status.Error(codes.NotFound, "URL not found")

		default:
			h.URLHandler.Logger.Error("Storage error: " + err.Error())
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	return &pb.ExpandResponse{OriginalUrl: &originalURL}, nil
}
