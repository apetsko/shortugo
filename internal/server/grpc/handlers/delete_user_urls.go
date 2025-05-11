package handlers

import (
	"context"

	"github.com/apetsko/shortugo/internal/models"
	pb "github.com/apetsko/shortugo/proto"
)

// DeleteUserURLs handles a request to delete multiple shortened URLs for a specific user.
// It asynchronously sends the deletion request to the internal processing channel.
//
// This method corresponds to the HTTP DELETE /api/user/urls endpoint.
//
// Request:
//   - user_id: string identifying the user
//   - short_url_ids: list of short URL identifiers to delete
//
// Response:
//   - success: true if the request was accepted for processing
func (h *Handler) DeleteUserURLs(ctx context.Context, req *pb.DeleteUserURLsRequest) (*pb.DeleteUserURLsResponse, error) {
	go func() {
		h.URLHandler.ToDelete <- models.BatchDeleteRequest{
			Ids:    req.ShortUrlIds,
			UserID: req.UserId,
		}
	}()

	return &pb.DeleteUserURLsResponse{Success: true}, nil
}
