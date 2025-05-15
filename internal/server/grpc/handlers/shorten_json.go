package handlers

import (
	"context"
	"errors"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenJSON creates a short URL from a single original URL.
// If the URL is already stored, it returns Conflict. Otherwise, it stores and returns the new short URL.
func (h *Handler) ShortenJSON(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetOriginalUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "original_url is required")
	}
	idLen := 8
	id := utils.GenerateID(req.GetOriginalUrl(), idLen)
	shortURL := h.URLHandler.BaseURL + "/" + id

	record := models.URLRecord{
		ID:     id,
		URL:    req.GetOriginalUrl(),
		UserID: req.GetUserId(),
	}

	// Check if it already exists
	existing, err := h.URLHandler.Storage.Get(ctx, id)
	if err == nil && existing != "" {
		return &pb.ShortenResponse{
			ShortUrl: &shortURL,
		}, status.Error(codes.AlreadyExists, "URL already exists")
	}

	if err != nil && !errors.Is(err, shared.ErrNotFound) {
		h.URLHandler.Logger.Error("failed to check existence", "error", err.Error())
		return nil, status.Error(codes.Internal, "failed to check URL")
	}

	// Store new URL
	if err := h.URLHandler.Storage.Put(ctx, record); err != nil {
		h.URLHandler.Logger.Error("failed to store URL", "error", err.Error())
		return nil, status.Error(codes.Internal, "failed to store URL")
	}

	return &pb.ShortenResponse{
		ShortUrl: &shortURL,
	}, nil
}
