package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Shorten accepts a raw URL string and returns a shortened version.
// If the URL already exists, it returns the same short URL with a Conflict code.
func (h *Handler) Shorten(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetOriginalUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "original_url is required")
	}
	idLen := 8
	id := utils.GenerateID(req.GetOriginalUrl(), idLen)
	record := models.URLRecord{
		ID:     id,
		URL:    req.GetOriginalUrl(),
		UserID: req.GetUserId(),
	}
	shortURL := h.URLHandler.BaseURL + "/" + id

	existing, err := h.URLHandler.Storage.Get(ctx, id)
	if err == nil && existing != "" {
		// Already exists
		fmt.Println("Already exists", existing)
		return &pb.ShortenResponse{
			ShortUrl: &shortURL,
		}, status.Error(codes.AlreadyExists, "URL already exists")
	}

	if err != nil && !errors.Is(err, shared.ErrNotFound) {
		h.URLHandler.Logger.Error("Get failed", "error", err.Error())
		return nil, status.Error(codes.Internal, "failed to get URL")
	}

	if err := h.URLHandler.Storage.Put(ctx, record); err != nil {
		h.URLHandler.Logger.Error("Put failed", "error", err.Error())
		return nil, status.Error(codes.Internal, "failed to store URL")
	}

	return &pb.ShortenResponse{
		ShortUrl: &shortURL,
	}, nil
}
