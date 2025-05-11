package handlers

import (
	"context"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/utils"
	pb "github.com/apetsko/shortugo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenBatch handles batch URL shortening requests.
// It validates and stores each original URL, and returns their shortened versions with correlation IDs.
func (h *Handler) ShortenBatch(ctx context.Context, req *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	var records []models.URLRecord
	var results []*pb.URLPair

	for _, item := range req.Urls {
		if item.OriginalUrl == "" {
			results = append(results, &pb.URLPair{
				CorrelationId: item.CorrelationId,
				ShortUrl:      "Bad Request: Empty URL",
			})
			continue
		}
		idLen := 8
		id := utils.GenerateID(item.OriginalUrl, idLen)

		record := models.URLRecord{
			URL:    item.OriginalUrl,
			ID:     id,
			UserID: req.UserId,
		}
		records = append(records, record)

		results = append(results, &pb.URLPair{
			CorrelationId: item.CorrelationId,
			ShortUrl:      h.URLHandler.BaseURL + "/" + id,
			OriginalUrl:   item.OriginalUrl, // optionally for debug
		})
	}

	if err := h.URLHandler.Storage.PutBatch(ctx, records); err != nil {
		h.URLHandler.Logger.Error("failed to store batch", "error", err.Error())
		return nil, status.Error(codes.Internal, "failed to store URLs")
	}

	return &pb.ShortenBatchResponse{
		Results: results,
	}, nil
}
