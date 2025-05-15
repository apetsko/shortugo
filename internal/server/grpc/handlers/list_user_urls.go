package handlers

import (
	"context"
	"errors"

	"github.com/apetsko/shortugo/internal/storages/shared"
	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListUserURLs returns all shortened URLs associated with the given user ID.
//
// Request:
//   - user_id: string
//
// Response:
//   - repeated URLPair (short + original URLs)
func (h *Handler) ListUserURLs(ctx context.Context, req *pb.ListUserURLsRequest) (*pb.ListUserURLsResponse, error) {
	records, err := h.URLHandler.Storage.ListLinksByUserID(ctx, h.URLHandler.BaseURL, req.GetUserId())
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			h.URLHandler.Logger.Error("no URLs for user: " + req.GetUserId())
			return nil, status.Error(codes.NotFound, "no URLs found for user")
		}
		h.URLHandler.Logger.Error("storage error: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to list URLs")
	}

	resp := &pb.ListUserURLsResponse{}
	for _, record := range records {
		empty := ""
		resp.Urls = append(resp.Urls, &pb.URLPair{
			CorrelationId: &empty, // not used here
			OriginalUrl:   &record.URL,
			ShortUrl:      &record.ID,
		})
	}

	return resp, nil
}
